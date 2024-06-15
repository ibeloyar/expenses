package password

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyPass = errors.New("password cannot be an empty string")
)

type PasswordManager struct {
	passCost int
}

func New(passCost int) *PasswordManager {
	return &PasswordManager{
		passCost: passCost,
	}
}

func (pm *PasswordManager) HashPassword(password string) (string, error) {
	if len(password) < 1 {
		return "", ErrEmptyPass
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.passCost)
	return string(bytes), err
}

func (pm *PasswordManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
