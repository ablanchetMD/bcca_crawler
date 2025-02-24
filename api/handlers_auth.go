package api

import (
	"bcca_crawler/internal/auth"
	"bcca_crawler/internal/caching"
	"bcca_crawler/internal/config"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/json_utils"
	"bcca_crawler/internal/auth/roles"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID         uuid.UUID     `json:"id"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
	Email      string        `json:"email"`
	Role       string        `json:"role"`
	IsVerified bool          `json:"is_verified"`
	DeletedAt  sql.NullTime  `json:"deleted_at"`
	DeletedBy  uuid.NullUUID `json:"deleted_by"`
	LastActive sql.NullTime  `json:"last_active"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,passwordstrength"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUserRequest struct {
	Email      string        `json:"email" validate:"omitempty,email"`
	Password   string        `json:"password" validate:"omitempty,passwordstrength"`
	Role       string        `json:"role" validate:"omitempty,oneof=user admin"`
	IsVerified *bool         `json:"is_verified" validate:"omitempty,boolean"`
	DeletedAt  sql.NullTime  `json:"deleted_at"`
	DeletedBy  uuid.NullUUID `json:"deleted_by"`
	LastActive sql.NullTime  `json:"last_active"`
}

type RevokeRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type Users struct {
	Users []User `json:"users"`
}


func mapUserStruct(src database.User) User {
	return User{
		ID:         src.ID,
		CreatedAt:  src.CreatedAt,
		UpdatedAt:  src.UpdatedAt,
		Email:      src.Email,
		Role:       src.Role,
		IsVerified: src.IsVerified,
		DeletedAt:  src.DeletedAt,
		DeletedBy:  src.DeletedBy,
		LastActive: src.LastActive,
	}
}

func HandleCLICreateUser(c *config.Config, email string, password string) (User, error) {
	var requestData = CreateUserRequest{
		Email:    email,
		Password: password,
	}
	err := c.Validate.Struct(requestData)
	if err != nil {
		return User{}, fmt.Errorf("validation failed: %s", err.Error())
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
		Role:      "admin",
	})
	if err != nil {
		return User{}, fmt.Errorf("error creating user")
	}
	return mapUserStruct(user), nil
}

func HandleCreateUser(c *config.Config, w http.ResponseWriter, r *http.Request) {

	var requestData CreateUserRequest
	err := UnmarshalAndValidatePayload(c, r, &requestData)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(requestData.Password)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user, err := c.Db.CreateUser(r.Context(), database.CreateUserParams{
		Email:     requestData.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Password:  hashedPassword,
		Role:      "user",
	})
	if err != nil {
		fmt.Println("Error creating user: ", err)
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating user")
		return
	}
	json_utils.RespondWithJSON(w, http.StatusCreated, mapUserStruct(user))
}

func HandleUpdateUser(c *config.Config, w http.ResponseWriter, r *http.Request) {
	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req UpdateUserRequest
	err = UnmarshalAndValidatePayload(c, r, &req)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()
	// Helper function to add updates to the map
	addUpdate := func(key string, value interface{}) {
		if value != nil && value != "" && !(value == (sql.NullTime{}) || value == (uuid.NullUUID{})) {
			updates[key] = value
		}
	}

	// Apply conditions for each field
	addUpdate("email", req.Email)
	if req.Password != "" {
		hashedPassword, err := auth.HashPassword(req.Password)
		if err != nil {
			json_utils.RespondWithError(w, http.StatusInternalServerError, "Error hashing password")
			return
		}
		addUpdate("password", hashedPassword)
	}
	addUpdate("role", req.Role)
	if req.IsVerified != nil {
		addUpdate("is_verified", req.IsVerified)
	}
	addUpdate("deleted_at", req.DeletedAt)
	addUpdate("deleted_by", req.DeletedBy)
	addUpdate("last_active", req.LastActive)

	user, err := UpdateUserDynamic(c, parsed_id.ID, updates, r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error updating user: %s with error:%s", parsed_id.ID.String(), err.Error()))
		return
	}

	json_utils.RespondWithJSON(w, http.StatusOK, user)
}

func HandleReset(c *config.Config, w http.ResponseWriter, r *http.Request) {
	if c.Platform != "dev" {
		json_utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to use this function.")
		return
	}

	err := c.Db.DeleteUsers(r.Context())
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting users")
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, "Users deleted")
}

//Refresh function uses cookies instead of headers, can we make it use both?

func HandleRefresh(c *config.Config, w http.ResponseWriter, r *http.Request) {	

	refresh_value, err := auth.GetCookieToken(r,"refresh_token")
	if err != nil {			
		json_utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token")
		return
	}
	

	usertoken, err := auth.ValidateRefreshToken(refresh_value, c)
	
	if err != nil {
		json_utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Refresh Token")
		return
	}

	auth.SetAuthCookies(w, usertoken)

	json_utils.RespondWithJSON(w, http.StatusOK, usertoken)
}

//Finish the Revoke Function

func HandleRevoke(c *config.Config, w http.ResponseWriter, r *http.Request) {
	var requestData RevokeRequest
	err := UnmarshalAndValidatePayload(c, r, &requestData)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = c.Db.RevokeRefreshTokenByUserId(r.Context(), requestData.UserID)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error revoking token")
		return
	}
	caching.DeleteRoleCache(requestData.UserID)
	json_utils.RespondWithJSON(w, http.StatusNoContent, "")
}

func HandleGetUsers(c *config.Config, q QueryParams, w http.ResponseWriter, r *http.Request) {	
	user,err := auth.GetUserFromContext(r)
	if err != nil {		
		json_utils.RespondWithError(w, http.StatusForbidden, "Invalid token")
		return
	}
	if user.Role != roles.Admin {
		json_utils.RespondWithError(w, http.StatusForbidden, "You are not authorized to use this function.")
		return
	}

	var users = []database.User{}
	params := database.GetUsersParams{
		Limit:  int32(q.Limit),
		Offset: int32(q.Offset),
	}
	
	//optional queries : sort, sort_by, page, limit, offset, filter, fields, include, exclude,
	switch {
	case q.FilterBy == "role" && len(q.Include) > 0:
		users, err = c.Db.GetUsersByRole(r.Context(), database.GetUsersByRoleParams{
			Role:   q.FilterBy,
			Limit:  params.Limit,
			Offset: params.Offset,
		})
	default:
		users, err = c.Db.GetUsers(r.Context(), params)
	}

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error getting users")
		return
	}
	var userStructs []User
	for _, user := range users {
		userStructs = append(userStructs, mapUserStruct(user))
	}
	json_utils.RespondWithJSON(w, http.StatusOK, Users{Users: userStructs})
}

func HandleLogin(c *config.Config, w http.ResponseWriter, r *http.Request) {
	var requestData LoginRequest
	err := UnmarshalAndValidatePayload(c, r, &requestData)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	
	user, err := c.Db.GetUserByEmail(r.Context(), requestData.Email)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusUnauthorized, "Password or email is invalid.")
		return
	}
	err = auth.CheckPasswordHash(requestData.Password, user.Password)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusUnauthorized, "Password or email is invalid.")
		return
	}

	refresh_token, err := auth.MakeRefreshToken()
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	refresh_obj, err := c.Db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating refresh token")
		return
	}

	user_tokens,err := auth.GetJWTFromRefreshToken(refresh_obj, c)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error creating auth token")
		return
	}
	user_role,err := roles.RoleFromString(user.Role)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Invalid Role")
		return
	}

	user_tokens.Role = user_role
	caching.SetRoleCache(user_tokens.UserID, user_role, time.Now().Add(time.Minute*60))

	auth.SetAuthCookies(w, user_tokens)
	json_utils.RespondWithJSON(w, http.StatusOK, user_tokens)
}

func HandleDeleteUserById(c *config.Config, w http.ResponseWriter, r *http.Request) {
	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = c.Db.DeleteUserByID(r.Context(), parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting user")
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": fmt.Sprintf("User %s deleted", parsed_id.ID.String())})
}

func HandleGetUserById(c *config.Config, w http.ResponseWriter, r *http.Request) {

	parsed_id, err := ParseAndValidateID(r)
	if err != nil {
		json_utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.Db.GetUserByID(r.Context(), parsed_id.ID)

	if err != nil {
		json_utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting user: %s", parsed_id.ID.String()))
		return
	}
	json_utils.RespondWithJSON(w, http.StatusOK, mapUserStruct(user))
}

func UpdateUserDynamic(c *config.Config, userID uuid.UUID, updates map[string]interface{}, r *http.Request) (User, error) {
	ctx := r.Context()
	DB := c.Database
	query := "UPDATE users SET "
	args := []interface{}{}
	i := 1

	for column, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", column, i)
		args = append(args, value)
		i++
	}

	query = strings.TrimSuffix(query, ", ") // Remove trailing comma
	query += fmt.Sprintf(" WHERE id = $%d ", i)
	args = append(args, userID)
	query += "RETURNING *;"
	fmt.Println(query)
	fmt.Println("args: \n", args)

	// user, err := db.Exec(query, args...)

	row := DB.QueryRowContext(ctx, query, args...)
	var password string
	var u User
	err := row.Scan(
		&u.ID,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.Email,
		&password, // password
		&u.Role,
		&u.IsVerified,
		&u.DeletedAt,
		&u.DeletedBy,
		&u.LastActive,
	)

	return u, err

}
