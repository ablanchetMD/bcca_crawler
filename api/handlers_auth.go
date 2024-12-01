package api

import (
	
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/auth"
	"github.com/google/uuid"
	"fmt"	
	"net/http"
	"context"
	"time"
)

type User struct {
	ID         uuid.UUID `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Email string    `json:"tumor_group"`	
}

type CreateUserRequest struct {
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password" validate:"required,passwordstrength"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`	
}

type Users struct {
	Users []User `json:"users"`
}

func mapUserStruct(src database.User) User {
	return User{
		ID:         src.ID,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		Email: 		src.Email,		
	}
}

func HandleCLICreateUser(c *config.Config, email string, password string) (User, error) {
	var requestData = CreateUserRequest{
		Email:    email,
		Password: password,
	}
	err := c.Validate.Struct(requestData)
	if err != nil {
		return User{}, fmt.Errorf("Validation failed: "+err.Error())
	}

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return User{}, fmt.Errorf("error hashing password")
	}
	user, err := c.Db.CreateUser(context.Background(), database.CreateUserParams{
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  hashedPassword,
		Role: "admin",
	})
	if err != nil {
		return User{}, fmt.Errorf("error creating user")
	}
	return mapUserStruct(user), nil
}

func HandleCreateUser(c *config.Config, w http.ResponseWriter, r *http.Request) {
	
	var requestData CreateUserRequest
	err :=UnmarshalAndValidatePayload(c, r, &requestData)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(requestData.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}
	
	user, err := c.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email:     requestData.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  hashedPassword,
		Role: "user",		
	})
	if err != nil {
		fmt.Println("Error creating user: ", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	respondWithJSON(w, http.StatusCreated, mapUserStruct(user))
}

func HandleReset(c *config.Config, w http.ResponseWriter, r *http.Request) {
	if c.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "You are not authorized to use this function.")
		return
	}

	err := c.Db.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users")
		return
	}
	respondWithJSON(w, http.StatusOK, "Users deleted")
}

func HandleRefresh(c *config.Config, w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	refreshToken, err := c.Db.GetRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked")
		return
	}

	jwt, err := auth.MakeJWT(refreshToken.UserID, c.Secret, time.Second*3600)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}
	//revoke old refresh token and issue new one, for rotation.
	_, err = c.Db.RevokeRefreshToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}
		_, err = c.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    refreshToken.UserID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
		// 60 days
		MaxAge: 5184000,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    jwt,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
		MaxAge:   3600, // 1 hour
	})

	t := struct {Token string `json:"auth_token"`;RefreshToken string `json:"refresh_token"`}{Token: jwt, RefreshToken: refresh_token}

	respondWithJSON(w, http.StatusOK, t)
}

func HandleRevoke(c *config.Config, w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token provided")
		return
	}

	_, err = c.Db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}
	respondWithJSON(w, http.StatusNoContent, "")
}

func HandleGetUsers(c *config.Config, w http.ResponseWriter, r *http.Request) {
	users, err := c.Db.GetUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}
	var userStructs []User
	for _, user := range users {
		userStructs = append(userStructs, mapUserStruct(user))
	}
	respondWithJSON(w, http.StatusOK, Users{Users: userStructs})
}

func HandleLogin(c *config.Config, w http.ResponseWriter, r *http.Request) {
	var requestData LoginRequest
	err :=UnmarshalAndValidatePayload(c, r, &requestData)	
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	expiresIn := time.Second * 3600

	user, err := c.Db.GetUserByEmail(r.Context(), requestData.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password or email is invalid.")
		return
	}
	err = auth.CheckPasswordHash(requestData.Password, user.Password)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Password or email is invalid.")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	_, err = c.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	token, err := auth.MakeJWT(user.ID, c.Secret, expiresIn)
	// user.Password = nil
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh_token,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
		// 60 days
		MaxAge: 5184000,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
		MaxAge:   3600, // 1 hour
	})

	t := struct {Token string `json:"auth_token"`;RefreshToken string `json:"refresh_token"`}{Token: token, RefreshToken: refresh_token}
	

	respondWithJSON(w, http.StatusOK, t)

}