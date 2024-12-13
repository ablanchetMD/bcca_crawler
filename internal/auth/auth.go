package auth

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"bcca_crawler/internal/config"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"bcca_crawler/internal/database"
	"bcca_crawler/internal/caching"
	"bcca_crawler/internal/auth/roles"	
	"encoding/hex"
	"context"
)

type contextKey string // Define your own type

const UserIDKey contextKey = "userID"     // Use a constant of your custom type
const UserRoleKey contextKey = "userRole" // Use a constant of your custom type

type UserToken struct {
	UserID         	uuid.UUID  `json:"id"`
	RefreshToken  	string     `json:"user_token"`
	AuthToken  		string     `json:"auth_token"`
	Role			roles.Role	   `json:"role"`
}



func HashPassword(password string) (string, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("HashPassword Function: %w",err)
	}
	return string(hashed_password),nil
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GetUserFromContext(r *http.Request) (UserToken, error) {
	ctx := r.Context()
	fmt.Println("Context: ",ctx)
	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		return UserToken{}, fmt.Errorf("GetUserFromContext Function: %s","user_id not found in context")
	}
	userRole, ok := r.Context().Value(UserRoleKey).(roles.Role)
	if !ok {
		return UserToken{}, fmt.Errorf("GetUserFromContext Function: %s","user_role not found in context")
	}
	user := UserToken{
		UserID: userID,
		Role: userRole,
	}
	return user, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error){
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "leukosys-auth",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	})
	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("MakeJWT Function: %w",err)
	}

	return tokenString, nil

}


func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token,err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("ValidateJWT Function: %w",err)
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return uuid.Nil, fmt.Errorf("ValidateJWT Function: %w",err)
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("ValidateJWT Function: %w",err)
	}

	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no Authorization header provided")
	}
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", fmt.Errorf("invalid Authorization header format")
	}
	return strings.TrimPrefix(authHeader, bearerPrefix), nil
}

func GetCookieToken(r *http.Request,token string) (string, error) {
	cookie, err := r.Cookie(token)
	if err != nil {
		return "", fmt.Errorf("GetCookieToken Function: %w",err)
	}
	return cookie.Value, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		fmt.Println("error with :",authHeader)
		return "", fmt.Errorf("no authorization header provided")
	}

	const apiPrefix = "ApiKey "
	if !strings.HasPrefix(authHeader, apiPrefix) {
		fmt.Println("error with :",authHeader,apiPrefix)
		return "", fmt.Errorf("invalid Authorization header format")
	}	
	return strings.TrimPrefix(authHeader, apiPrefix), nil	
}

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	_,err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("MakeRefreshToken Function: %w",err)
	}
	return hex.EncodeToString(randomBytes), nil
}

func AddCookie(w http.ResponseWriter, name, value string, timeInSeconds int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   true, // Set to false if not using HTTPS
		Path:     "/",
		MaxAge:   timeInSeconds,
	})
}

func SetAuthCookies(w http.ResponseWriter, user UserToken) {
	AddCookie(w, "refresh_token", user.RefreshToken, 5184000)
	AddCookie(w, "auth_token", user.AuthToken, 120)
}


func GetJWTFromRefreshToken(token database.RefreshToken, c *config.Config) (UserToken, error) {	
	user := UserToken{}	

	jwt, err := MakeJWT(token.UserID, c.Secret, time.Second*120)
	if err != nil {
		return user, fmt.Errorf("GetJWTFromRefreshToken Function : %w",err)
	}	

	user.RefreshToken = token.Token
	user.AuthToken = jwt
	user.UserID = token.UserID

	return user, nil
}


func ValidateRefreshToken(token string, c *config.Config) (UserToken, error) {
	ctx := context.Background()
	user := UserToken{}

	refreshToken, err := c.Db.GetRefreshToken(ctx, token)
	if err != nil {		
		return user, fmt.Errorf("ValidateRefreshToken Function : %w",err)
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		return user, fmt.Errorf("ValidateRefreshToken Function : %s","token expired")
	}

	if refreshToken.RevokedAt.Valid {
		return user, fmt.Errorf("ValidateRefreshToken Function : %s","token revoked")
	}	

	jwt, err := MakeJWT(refreshToken.UserID, c.Secret, time.Second*120)
	if err != nil {
		return user, fmt.Errorf("ValidateRefreshToken Function : %w",err)
	}
	//revoke old refresh token and issue new one, for rotation.
	_, err = c.Db.RevokeRefreshToken(ctx, refreshToken.Token)
	if err != nil {
		return user, fmt.Errorf("ValidateRefreshToken Function : %w",err)
	}

	refresh_token, err := MakeRefreshToken()
	if err != nil {
		return user, fmt.Errorf("ValidateRefreshToken Function : %w",err)
	}
	_, err = c.Db.CreateRefreshToken(ctx, database.CreateRefreshTokenParams{
		Token:     refresh_token,
		UserID:    refreshToken.UserID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	})

	if err != nil {
		return user, fmt.Errorf("ValidateRefreshToken Function : %w",err)
	}

	userRole,err := roles.RoleFromString(refreshToken.Role)
	if err != nil {
		return user, fmt.Errorf("ValidateRefreshToken Function: %w",err)
	}

	user.RefreshToken = refresh_token
	user.AuthToken = jwt
	user.UserID = refreshToken.UserID
	user.Role = userRole
	caching.SetRoleCache(refreshToken.UserID, userRole, time.Now().Add(time.Minute * 60))

	return user, nil
}