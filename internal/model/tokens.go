package model

import "time"

type Token struct {
	UserID    int        `json:"userID"`
	Token     string     `json:"token"`
	CreatedAt *time.Time `json:"createdAt"`
}
