package model

import "time"

type User struct {
	ID             int        `json:"id"`
	Login          string     `json:"login"`
	Email          string     `json:"email"`
	EmailConfirmed bool       `json:"emailConfirmed"`
	Password       string     `json:"password"`
	RoleID         int        `json:"roleID"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

type CreateUserBody struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EditUserBody struct {
	Login string `json:"login"`
	Email string `json:"email"`
}

type EditUserPasswordBody struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
}
