package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"message/internal/storage"
	"results/errs"
)

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
