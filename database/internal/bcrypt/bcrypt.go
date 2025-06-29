package bcrypt

import (
	"database/internal/db_models"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"results/errs"
)

func Encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func ValidatePassword(user db_models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return errs.WrongPassword
		}
		return err
	}

	return nil
}
