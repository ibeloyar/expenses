package http

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/services/auth"
	"github.com/B-Dmitriy/expenses/internal/services/categories"
	"github.com/B-Dmitriy/expenses/internal/services/counterparties"
	"github.com/B-Dmitriy/expenses/internal/services/mail"
	"github.com/B-Dmitriy/expenses/internal/services/transactions"
	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/tokens"
	"github.com/go-playground/validator/v10"

	categoriesDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/categories"
	counterpartiesDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/counterparties"
	tokensDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/tokens"
	transactionsDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/transactions"
	usersDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/users"
)

func initRoutes(
	cfg *config.Config,
	serv *http.ServeMux,
	logger *slog.Logger,
	db *postgres.PGStorage,
	tm *tokens.TokensManager,
	pm *password.PasswordManager,
) *http.ServeMux {
	v := validator.New()
	utils := postgres.NewPGUtils()

	usersStore := usersDB.NewUsersStorage(db)
	tokensStore := tokensDB.NewTokensStorage(db)
	categoriesStore := categoriesDB.NewCategoriesStorage(db)
	transactionsStore := transactionsDB.NewTransactionsStorage(db)
	counterpartiesStore := counterpartiesDB.NewCounterpartiesStorage(db)

	usersService := users.NewUsersService(logger, usersStore, v, utils, pm)
	mailService := mail.NewMailService(logger, &cfg.Mail, &cfg.HTTPServer, usersStore)
	authService := auth.NewAuthService(logger, utils, v, usersStore, tokensStore, tm, pm)
	categoriesService := categories.NewCategoriesService(logger, categoriesStore, v, utils)
	transactionsService := transactions.NewTransactionsService(logger, transactionsStore, v, utils)
	counterpartiesService := counterparties.NewCounterpartiesService(logger, counterpartiesStore, v, utils)

	// CORS
	serv.Handle("OPTIONS /*", CorsMiddleware(http.HandlerFunc(CorsOptionHandlerFunc)))

	// Auth
	serv.Handle(
		"POST /api/v1/login",
		CorsMiddleware(http.HandlerFunc(authService.Login)),
	)
	serv.Handle(
		"POST /api/v1/registration",
		CorsMiddleware(http.HandlerFunc(authService.Registration)),
	)
	serv.Handle(
		"POST /api/v1/logout",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(authService.Logout))),
	)
	serv.Handle(
		"POST /api/v1/refresh",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(authService.Refresh))),
	)

	// Users
	serv.Handle(
		"GET /api/v1/users/{userID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(usersService.GetUser))),
	)
	serv.Handle(
		"PUT /api/v1/users/{userID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(usersService.EditUserInfo))),
	)
	serv.Handle(
		"DELETE /api/v1/users/{userID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(usersService.DeleteUser))),
	)
	serv.Handle(
		"GET /api/v1/users",
		CorsMiddleware(authService.AuthOnlyAdminMiddleware(http.HandlerFunc(usersService.GetUsersList))),
	)
	serv.Handle(
		"POST /api/v1/users",
		CorsMiddleware(authService.AuthOnlyAdminMiddleware(http.HandlerFunc(usersService.CreateUser))),
	)

	// Categories
	serv.Handle(
		"GET /api/v1/categories",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(categoriesService.GetCategoriesList))),
	)
	serv.Handle(
		"GET /api/v1/categories/{categoryID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(categoriesService.GetCategoryByID))),
	)
	serv.Handle(
		"POST /api/v1/categories",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(categoriesService.CreateCategory))),
	)
	serv.Handle(
		"PUT /api/v1/categories/{categoryID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(categoriesService.EditCategory))),
	)
	serv.Handle(
		"DELETE /api/v1/categories/{categoryID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(categoriesService.DeleteCategory))),
	)

	// Counterparties
	serv.Handle(
		"GET /api/v1/counterparties",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(counterpartiesService.GetCounterpartiesList))),
	)
	serv.Handle(
		"GET /api/v1/counterparties/{counterpartyID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(counterpartiesService.GetCounterpartyByID))),
	)
	serv.Handle(
		"POST /api/v1/counterparties",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(counterpartiesService.CreateCounterparty))),
	)
	serv.Handle(
		"PUT /api/v1/counterparties/{counterpartyID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(counterpartiesService.EditCounterparty))),
	)
	serv.Handle(
		"DELETE /api/v1/counterparties/{counterpartyID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(counterpartiesService.DeleteCounterparty))),
	)

	// Transactions
	serv.Handle(
		"GET /api/v1/transactions",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(transactionsService.GetTransactionsList))),
	)
	serv.Handle(
		"GET /api/v1/transactions/{transactionID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(transactionsService.GetTransactionByID))),
	)
	serv.Handle(
		"POST /api/v1/transactions",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(transactionsService.CreateTransaction))),
	)
	serv.Handle(
		"PUT /api/v1/transactions/{transactionID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(transactionsService.EditTransaction))),
	)
	serv.Handle(
		"DELETE /api/v1/transactions/{transactionID}",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(transactionsService.DeleteTransaction))),
	)

	// Mail
	serv.Handle(
		"GET /api/v1/confirm:send",
		CorsMiddleware(authService.AuthMiddleware(http.HandlerFunc(mailService.RequestConfirmMail))),
	)
	serv.Handle(
		"GET /api/v1/confirm:approve",
		CorsMiddleware(http.HandlerFunc(mailService.ConfirmUserAccount)),
	)

	return serv
}
