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
	ConfirmToken   *string    `json:"confirmToken"`
}

type UserInfo struct {
	ID             int        `json:"id"`
	Login          string     `json:"login"`
	Email          string     `json:"email"`
	EmailConfirmed bool       `json:"emailConfirmed"`
	RoleID         int        `json:"roleID"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

type CreateUserBody struct {
	Login    string `json:"login" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=4,max=20"`
}

type EditUserBody struct {
	Login string `json:"login" validate:"required,min=2,max=255"`
	Email string `json:"email" validate:"required,email"`
}

type EditUserPasswordBody struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=4,max=20"`
	NewPassword string `json:"newPassword" validate:"required,min=4,max=20"`
}
