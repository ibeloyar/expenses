package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/B-Dmitriy/expenses/internal/storage/mysql"
	"github.com/gorilla/mux"
)

type UsersHandler struct {
	storage *mysql.Storage
	logger  *slog.Logger
}

func NewUsersHandler(storage *mysql.Storage, logger *slog.Logger) *UsersHandler {
	return &UsersHandler{
		storage: storage,
		logger:  logger,
	}
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	const os = "server.http.users.GetUsers"

	users, err := h.storage.GetUsers(5, 0)
	if err != nil {
		h.logger.Error("%s: %w", os, err)
	}

	res, err := json.Marshal(users)
	if err != nil {
		h.logger.Error("%s: %w", os, err)
	}

	w.Write(res)
}

func (h *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	const os = "server.http.users.GetUser"
	id := mux.Vars(r)["id"]

	intID, err := strconv.Atoi(id)
	if err != nil {
		h.logger.Error("%s: %w", os, err)
	}

	user, err := h.storage.GetUser(intID)
	if err != nil {
		h.logger.Error("%s: %w", os, err)
	}

	res, err := json.Marshal(user)
	if err != nil {
		h.logger.Error("%s: %w", os, err)
	}

	w.Write(res)
}
