package http

import (
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/services/users"
	"github.com/jackc/pgx/v5"

	usersDB "github.com/B-Dmitriy/expenses/internal/storage/users"
)

func initRoutes(serv *http.ServeMux, l *slog.Logger, db *pgx.Conn) *http.ServeMux {
	usersStore := usersDB.NewUsersStorage(db)

	usersService := users.NewUsersService(l, usersStore)

	serv.HandleFunc("GET /api/v1/users", usersService.GetUsersList)
	serv.HandleFunc("POST /api/v1/users", usersService.CreateUser)
	serv.HandleFunc("GET /api/v1/users/{userID}", usersService.GetUser)
	serv.HandleFunc("PUT /api/v1/users/{userID}", usersService.EditUserInfo)
	serv.HandleFunc("DELETE /api/v1/users/{userID}", usersService.DeleteUser)

	return serv
}
