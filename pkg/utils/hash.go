package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
	ErrPasswordTooLong  = errors.New("password must be less than 72 characters")
)

// HashPassword generates a bcrypt hash from the given password
func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", ErrPasswordTooShort
	}
	if len(password) > 72 {
		return "", ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword compares a password with a hash
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
