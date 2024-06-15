package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/jackc/pgx/v5"

	us "github.com/B-Dmitriy/expenses/internal/storage/users"
)

type UsersService struct {
	logger      *slog.Logger
	store       *us.UsersStorage
	passManager *password.PasswordManager
}

func NewUsersService(
	l *slog.Logger,
	us *us.UsersStorage,
	pm *password.PasswordManager,
) *UsersService {
	return &UsersService{
		logger:      l,
		store:       us,
		passManager: pm,
	}
}

func (us *UsersService) GetUsersList(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, us.logger, "users.GetUsersList")

	p, err := web.ParseQueryPagination(r, &web.Pagination{Page: 1, Limit: 25})
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	search, err := web.ParseSearchString(r)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	users, err := us.store.GetUsersList(p.Page, p.Limit, search)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
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
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	user, err := us.store.GetUser(userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	web.WriteOK(w, user)
}

func (us *UsersService) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := new(model.CreateUserBody)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	passHash, err := us.passManager.HashPassword(body.Password)
	if err != nil {
		if errors.As(err, &password.ErrEmptyPass) {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	err = us.store.CreateUser(replacePasswordOnHash(body, passHash))
	if err != nil {
		if isConstrain, e := us.store.CheckConstrainErr(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
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
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	body := new(model.EditUserBody)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	err = us.store.EditUser(userID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		if isConstrain, err := us.store.CheckConstrainErr(err); isConstrain {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

func (us *UsersService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	err = us.store.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
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
