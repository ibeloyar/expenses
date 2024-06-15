package http

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/password"

	usersDB "github.com/B-Dmitriy/expenses/internal/storage/users"
)

func initRoutes(
	serv *http.ServeMux,
	l *slog.Logger,
	db *storage.PGStorage,
	pm *password.PasswordManager,
) *http.ServeMux {
	pgServiceUtils := storage.NewPGServiceUtils()
	usersStore := usersDB.NewUsersStorage(db)

	usersService := users.NewUsersService(l, usersStore, pgServiceUtils, pm)

	serv.HandleFunc("GET /api/v1/users", usersService.GetUsersList)
	serv.HandleFunc("POST /api/v1/users", usersService.CreateUser)
	serv.HandleFunc("GET /api/v1/users/{userID}", usersService.GetUser)
	serv.HandleFunc("PUT /api/v1/users/{userID}", usersService.EditUserInfo)
	serv.HandleFunc("DELETE /api/v1/users/{userID}", usersService.DeleteUser)

	return serv
}
