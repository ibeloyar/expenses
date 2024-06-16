package http

import (
	"github.com/B-Dmitriy/expenses/pgk/tokens"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/services/auth"
	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/go-playground/validator/v10"

	tokensDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/tokens"
	usersDB "github.com/B-Dmitriy/expenses/internal/storage/postgres/users"
)

func initRoutes(
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

	usersService := users.NewUsersService(logger, usersStore, v, utils, pm)
	authService := auth.NewAuthService(logger, utils, v, usersStore, tokensStore, tm, pm)

	serv.HandleFunc("POST /api/v1/login", authService.Login)
	serv.HandleFunc("POST /api/v1/logout", authService.Logout)
	serv.HandleFunc("POST /api/v1/refresh", authService.Refresh)
	serv.HandleFunc("POST /api/v1/registration", authService.Registration)

	// TODO: разобраться, какие данные может получать пользователь. И как с ними взаимодействовать.
	// Пользователь может: получить СВОЙ аккаунт, удалить СВОЙ аккаунт, изменить СВОЙ аккаунт
	// Создавать пользователя вручную, не через регистацию может только админ, так же как и просматривать всех юзеров
	serv.Handle("GET /api/v1/users", authService.AdminAuthMiddleware(http.HandlerFunc(usersService.GetUsersList)))
	serv.Handle("POST /api/v1/users", authService.AdminAuthMiddleware(http.HandlerFunc(usersService.CreateUser)))
	serv.HandleFunc("GET /api/v1/users/{userID}", usersService.GetUser)
	serv.HandleFunc("PUT /api/v1/users/{userID}", usersService.EditUserInfo)
	serv.HandleFunc("DELETE /api/v1/users/{userID}", usersService.DeleteUser)

	return serv
}
