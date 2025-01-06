package storage

import "github.com/ibeloyar/expenses/internal/model"

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

	AddConfirmToken(id int, confirmToken string) error
	ConfirmUserMail(confirmToken string) error
}

type TokensStore interface {
	GetTokenByUserID(userID int) (*model.Token, error)
	CheckToken(userID int) (bool, error)
	CreateToken(userID int, token string) error
	ChangeToken(userID int, token string) error
	DeleteToken(userID int) error
}

type CategoriesStore interface {
	GetAllUserCategories(userID, page, limit int, search string) ([]*model.Category, error)
	GetCategoryByID(id, userID int) (*model.Category, error)
	CreateCategory(userID int, data *model.CreateCategoryBody) error
	EditCategory(categoryID, userID int, data *model.EditCategoryBody) error
	DeleteCategory(id, userID int) error
}

type CounterpartiesStore interface {
	GetAllUserCounterparties(userID, page, limit int, search string) ([]*model.Counterparty, error)
	GetCounterpartyByID(id, userID int) (*model.Counterparty, error)
	CreateCounterparty(userID int, data *model.CreateCounterpartyBody) error
	EditCounterparty(counterpartyID, userID int, data *model.EditCounterpartyBody) error
	DeleteCounterparty(id, userID int) error
}

type TransactionsStore interface {
	GetAllUserTransactions(userID, page, limit int, search string) ([]*model.Transaction, error)
	GetTransactionByID(id, userID int) (*model.Transaction, error)
	CreateTransaction(userID int, data *model.CreateTransactionBody) error
	EditTransaction(transactionID, userID int, data *model.EditTransactionBody) error
	DeleteTransaction(id, userID int) error
}
