package model

import "time"

type Transaction struct {
	ID             int        `json:"id"`
	UserID         int        `json:"userID"`
	CategoryID     int        `json:"categoryID"`
	CounterpartyID int        `json:"counterpartyID"`
	Type           string     `json:"type"`
	Date           *time.Time `json:"date"`
	Amount         float64    `json:"amount"`
	Currency       string     `json:"currency"`
	Comment        string     `json:"comment"`
	CreatedAt      *time.Time `json:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt"`
}

type CreateTransactionBody struct {
	CategoryID     int        `json:"categoryID" validate:"required,min=1"`
	CounterpartyID int        `json:"counterpartyID" validate:"required,min=1"`
	Type           string     `json:"type" validate:"required,oneof=in out"`
	Date           *time.Time `json:"date" validate:"required"`
	Amount         float64    `json:"amount" validate:"required,min=0"`
	Currency       string     `json:"currency" validate:"required"`
	Comment        string     `json:"comment" validate:"max=2048"`
}

type EditTransactionBody struct {
	CategoryID     int        `json:"categoryID" validate:"required,min=1"`
	CounterpartyID int        `json:"counterpartyID" validate:"required,min=1"`
	Type           string     `json:"type" validate:"required,oneof=in out"`
	Date           *time.Time `json:"date" validate:"required"`
	Amount         float64    `json:"amount" validate:"required,min=0"`
	Currency       string     `json:"currency" validate:"required"`
	Comment        string     `json:"comment" validate:"max=2048"`
}
