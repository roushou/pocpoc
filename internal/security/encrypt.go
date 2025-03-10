package security

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password.
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty password")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), err
}
