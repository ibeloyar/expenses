package password

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
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

func (pm *PasswordManager) GetAuthorizationHeader(r *http.Request) (string, error) {
	headerToken := r.Header.Get("Authorization")
	bearerTokenSlice := strings.Split(headerToken, " ")

	if len(bearerTokenSlice) < 2 {
		return "", fmt.Errorf("bearer token not found")
	}

	return bearerTokenSlice[1], nil
}
