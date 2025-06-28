package jwt

import (
	"auth/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func GenerateToken() (string, error) {
	key := []byte(storage.Env.JwtSecret)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"Id":  "1",
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}
