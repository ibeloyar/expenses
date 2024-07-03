package model

import "time"

type Counterparty struct {
	ID          int        `json:"id"`
	UserID      int        `json:"userID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type CreateCounterpartyBody struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"max=1024"`
}

type EditCounterpartyBody struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description" validate:"max=1024"`
}
