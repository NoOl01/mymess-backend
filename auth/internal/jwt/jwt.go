package jwt

import (
	"auth/internal/storage"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"results/errs"
	"time"
)

func GenerateToken(userId string) (string, error) {
	key := []byte(storage.Env.JwtSecret)

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(30 * 24 * time.Hour).Unix(),
	})

	accessTokenString, err := accessToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func ValidateJwt(accessToken string) (string, error) {
	jwtSecret := []byte(storage.Env.JwtSecret)

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.UnexpectedSignMethod
		}
		return jwtSecret, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrInvalidKey) || !token.Valid {
			return "", errs.InvalidToken
		}
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errs.InvalidTokenClaimsType
	}

	userId, ok := claims["sub"].(string)
	if !ok {
		return "", fmt.Errorf("%w: sub", errs.InvalidOrMissingClaim)
	}

	return userId, nil
}
