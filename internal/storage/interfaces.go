package storage

import "github.com/B-Dmitriy/expenses/internal/model"

type UsersStore interface {
	GetUsersList() ([]*model.UserInfo, error)
	GetUser(id int) (*model.UserInfo, error)
	CreateUser(body *model.CreateUserBody) error
	EditUser(id int, user *model.EditUserBody) error
	DeleteUser(id int) error
}
