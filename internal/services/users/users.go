package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/jackc/pgx/v5"

	storage "github.com/B-Dmitriy/expenses/internal/storage/users"
)

type UsersService struct {
	logger      *slog.Logger
	store       *storage.UsersStorage
	passManager *password.PasswordManager
}

func NewUsersService(
	l *slog.Logger,
	s *storage.UsersStorage,
	pm *password.PasswordManager,
) *UsersService {
	return &UsersService{
		logger:      l,
		store:       s,
		passManager: pm,
	}
}

func (us *UsersService) GetUsersList(w http.ResponseWriter, r *http.Request) {
	users, err := us.store.GetUsersList()
	if err != nil {
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	web.WriteOK(w, users)
}

func (us *UsersService) GetUser(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	user, err := us.store.GetUser(userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	web.WriteOK(w, user)
}

func (us *UsersService) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := new(model.CreateUserBody)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	passHash, err := us.passManager.HashPassword(body.Password)
	if err != nil {
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	err = us.store.CreateUser(replacePasswordOnHash(body, passHash))
	if err != nil {
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	web.WriteOK(w, nil)
}

func (us *UsersService) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	body := new(model.EditUserBody)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	err = us.store.EditUser(userID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	web.WriteOK(w, nil)
}

func (us *UsersService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	err = us.store.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		us.logger.Error(err.Error())
		web.WriteServerError(w)
		return
	}

	web.WriteNoContent(w, nil)
}

func replacePasswordOnHash(user *model.CreateUserBody, hash string) *model.CreateUserBody {
	return &model.CreateUserBody{
		Login:    user.Login,
		Email:    user.Email,
		Password: hash,
	}
}
