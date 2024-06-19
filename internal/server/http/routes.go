package http

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/services/auth"
	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/tokens"
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

	// Auth
	serv.HandleFunc("POST /api/v1/login", authService.Login)
	serv.HandleFunc("POST /api/v1/registration", authService.Registration)
	serv.Handle("POST /api/v1/logout", authService.AuthMiddleware(http.HandlerFunc(authService.Logout)))
	serv.Handle("POST /api/v1/refresh", authService.AuthMiddleware(http.HandlerFunc(authService.Refresh)))

	// Users
	serv.Handle("GET /api/v1/users/{userID}", authService.AuthMiddleware(http.HandlerFunc(usersService.GetUser)))
	serv.Handle("PUT /api/v1/users/{userID}", authService.AuthMiddleware(http.HandlerFunc(usersService.EditUserInfo)))
	serv.Handle("DELETE /api/v1/users/{userID}", authService.AuthMiddleware(http.HandlerFunc(usersService.DeleteUser)))
	serv.Handle("GET /api/v1/users", authService.AuthOnlyAdminMiddleware(http.HandlerFunc(usersService.GetUsersList)))
	serv.Handle("POST /api/v1/users", authService.AuthOnlyAdminMiddleware(http.HandlerFunc(usersService.CreateUser)))

	return serv
}
