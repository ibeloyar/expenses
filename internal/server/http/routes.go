package http

import (
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/password"

	usersDB "github.com/B-Dmitriy/expenses/internal/storage/users"
)

func initRoutes(
	serv *http.ServeMux,
	logger *slog.Logger,
	db *storage.PGStorage,
	pm *password.PasswordManager,
) *http.ServeMux {
	v := validator.New()
	utils := storage.NewPGServiceUtils()

	usersStore := usersDB.NewUsersStorage(db)

	usersService := users.NewUsersService(logger, usersStore, v, utils, pm)

	serv.HandleFunc("GET /api/v1/users", usersService.GetUsersList)
	serv.HandleFunc("POST /api/v1/users", usersService.CreateUser)
	serv.HandleFunc("GET /api/v1/users/{userID}", usersService.GetUser)
	serv.HandleFunc("PUT /api/v1/users/{userID}", usersService.EditUserInfo)
	serv.HandleFunc("DELETE /api/v1/users/{userID}", usersService.DeleteUser)

	return serv
}
