package model

import "time"

type Category struct {
	ID          int        `json:"id"`
	UserID      int        `json:"userID"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

type CreateCategoryBody struct {
	UserID      int    `json:"userID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type EditCategoryBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
