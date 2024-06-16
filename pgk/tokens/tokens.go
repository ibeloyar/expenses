package tokens

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type Tokens struct {
	AcceptToken  string `json:"acceptToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokensManager struct {
	secretKey string
}

type UserInfo struct {
	UserID     int `json:"userID"`
	UserRoleID int `json:"userRoleID"`
}

func New(secretKey string) *TokensManager {
	return &TokensManager{
		secretKey: secretKey,
	}
}

func (tm *TokensManager) GenerateTokens(userID, userRoleID int) (*Tokens, error) {
	acceptTokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"userID":     userID,
		"userRoleID": userRoleID,
		"exp":        time.Now().Add(time.Hour * 2).Unix(), // 2 hours
	})

	refreshTokenData := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"userID":     userID,
		"userRoleID": userRoleID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // 1 day
	})

	acceptToken, err := acceptTokenData.SignedString([]byte(tm.secretKey))
	if err != nil {
		return nil, err
	}

	refreshToken, err := refreshTokenData.SignedString([]byte(tm.secretKey))
	if err != nil {
		return nil, err
	}

	return &Tokens{
		AcceptToken:  acceptToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tm *TokensManager) VerifyJWTToken(tokenString string) (*UserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return tm.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userID, err := strconv.Atoi(fmt.Sprintf("%v", claims["userID"]))
		if err != nil {
			return nil, err
		}
		userRoleID, err := strconv.Atoi(fmt.Sprintf("%v", claims["userRoleID"]))
		if err != nil {
			return nil, err
		}
		return &UserInfo{
			UserRoleID: userRoleID,
			UserID:     userID,
		}, nil
	} else {
		return nil, err
	}
}
