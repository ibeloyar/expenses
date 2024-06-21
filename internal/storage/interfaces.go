package storage

import "github.com/B-Dmitriy/expenses/internal/model"

type ServiceUtils interface {
	CheckConstrainError(e error) (bool, error)
}

type UsersStore interface {
	GetUsersList(page, limit int, search string) ([]*model.UserInfo, error)
	GetUser(id int) (*model.UserInfo, error)
	GetUserByEmail(email string) (*model.User, error)
	CreateUser(body *model.CreateUserBody) error
	EditUser(id int, user *model.EditUserBody) error
	DeleteUser(id int) error
}

type TokensStore interface {
	GetTokenByUserID(userID int) (*model.Token, error)
	CheckToken(userID int) (bool, error)
	CreateToken(userID int, token string) error
	ChangeToken(userID int, token string) error
	DeleteToken(userID int) error
}

type CategoriesStore interface {
	GetAllUserCategories(userID, page, limit int) ([]*model.Category, error)
	GetCategoryByID(id int) (*model.Category, error)
	CreateCategory(data *model.CreateCategoryBody) error
	EditCategory(categoryID int, data *model.EditCategoryBody) error
	DeleteCategory(id int) error
}
