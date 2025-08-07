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
		"type":   "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getKey()))
}

func SignRefreshToken(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
		"type":   "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(getKey()))
}

func VerifyPayload(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(getKey()), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(getKey()), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if tokenType, exists := claims["type"]; !exists || tokenType != "refresh" {
			return nil, jwt.ErrSignatureInvalid
		}
	}

	return token, nil
}
