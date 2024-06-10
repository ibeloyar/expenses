package storage

import "github.com/B-Dmitriy/expenses/internal/model"

type UsersStore interface {
	GetList() ([]model.User, error)
	GetUser(id int) (*model.User, error)
}
