package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/storage/mysql"
	"github.com/gorilla/mux"
)

func NewServer(cfg config.HTTPSettings, storage *mysql.Storage, logger *slog.Logger) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := mux.NewRouter()

	usersHandler := NewUsersHandler(storage, logger)

	r.HandleFunc("/api/v1/users", usersHandler.GetUsers)
	r.HandleFunc("/api/v1/users/{id}", usersHandler.GetUser)

	return &http.Server{
		Addr:         addr,
		Handler:      r,
		WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.IddleTimout) * time.Second,
	}
}
