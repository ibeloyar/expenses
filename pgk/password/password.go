package password

import (
	"golang.org/x/crypto/bcrypt"
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), pm.passCost)
	return string(bytes), err
}

func (pm *PasswordManager) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
