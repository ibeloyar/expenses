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
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

type UsersPGService struct {
	logger      *slog.Logger
	store       storage.UsersStore
	validator   *validator.Validate
	utils       storage.ServiceUtils
	passManager *password.PasswordManager
}

func NewUsersService(
	l *slog.Logger,
	us storage.UsersStore,
	v *validator.Validate,
	u storage.ServiceUtils,
	pm *password.PasswordManager,
) *UsersPGService {
	return &UsersPGService{
		logger:      l,
		store:       us,
		validator:   v,
		utils:       u,
		passManager: pm,
	}
}

func (us *UsersPGService) GetUsersList(w http.ResponseWriter, r *http.Request) {
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

func (us *UsersPGService) GetUser(w http.ResponseWriter, r *http.Request) {
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

func (us *UsersPGService) CreateUser(w http.ResponseWriter, r *http.Request) {
	body := new(model.CreateUserBody)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	err = us.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
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
		if isConstrain, e := us.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

func (us *UsersPGService) EditUserInfo(w http.ResponseWriter, r *http.Request) {
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

	err = us.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	err = us.store.EditUser(userID, body)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, fmt.Errorf("user %d not found", userID))
			return
		}
		if isConstrain, err := us.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	web.WriteOK(w, nil)
}

func (us *UsersPGService) DeleteUser(w http.ResponseWriter, r *http.Request) {
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
