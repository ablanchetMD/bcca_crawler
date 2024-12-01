package auth

import (
	"fmt"
	"net/http"
	"time"
	"strings"
	"crypto/rand"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"encoding/hex"
)

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