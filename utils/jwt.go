package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getKey() string {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		panic("jwt secret not found")
	}
	return key
}

func SignPayload(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(3 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getKey()))
}

func VerifyPayload(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(getKey()), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}
