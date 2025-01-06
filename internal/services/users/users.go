package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ibeloyar/expenses/internal/model"
	"github.com/ibeloyar/expenses/internal/storage"
	"github.com/ibeloyar/expenses/pgk/password"
	"github.com/ibeloyar/expenses/pgk/web"
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

const (
	AdminRoleID = 1 // TODO: вынести в конфиг
)

// GetUsersList
// @Router /api/v1/users [get]
// @Tags Users
// @Param page query int false "positive int" minimum(1) maximum(10) default(1)
// @Param limit query int false "positive int" minimum(1) maximum(100) default(25)
// @Param search query string false "any string" maxlength(256)
// @Description Получить список пользователей (только для админа)
// @Security BearerAuth
// @Success 200 {object} []model.UserInfo
// @Failure 400 {object} web.WebError
// @Failure 401 {object} web.WebError
// @Failure 403 {object} web.WebError
// @Failure 500 {object} web.WebError
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

// GetUser
// @Router /api/v1/users/{id} [get]
// @Tags Users
// @Param id path int true "User ID"
// @Description Получить информацию о пользователе (Пользователь - о семе, Админ о любом)
// @Security BearerAuth
// @Success 200 {object} model.UserInfo
// @Failure 400 {object} web.WebError
// @Failure 401 {object} web.WebError
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (us *UsersPGService) GetUser(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, us.logger, "users.GetUser")

	tokenUserID := r.Context().Value("userID")
	tokenUserRoleID := r.Context().Value("userRoleID")

	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	if tokenUserID != userID && tokenUserRoleID != AdminRoleID {
		web.WriteForbidden(w, nil)
		return
	}

	user, err := us.store.GetUser(userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	web.WriteOK(w, user)
}

// CreateUser
// @Router /api/v1/users [post]
// @Tags Users
// @Param request body model.CreateUserBody false "query params"
// @Description Создать пользователя вручную. (Только админ)
// @Security BearerAuth
// @Success 201
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (us *UsersPGService) CreateUser(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, us.logger, "users.CreateUser")

	body := new(model.CreateUserBody)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}
	defer r.Body.Close()

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

	web.WriteCreated(w, nil)
}

// EditUserInfo
// @Router /api/v1/users/{id} [put]
// @Tags Users
// @Param id path int true "User ID"
// @Param request body model.EditUserBody false "query params"
// @Description Изменить информацию о пользователе (Пользователь - о семе, Админ о любом)
// @Security BearerAuth
// @Success 200
// @Failure 400 {object} web.WebError
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (us *UsersPGService) EditUserInfo(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, us.logger, "users.EditUserInfo")

	tokenUserID := r.Context().Value("userID")
	tokenUserRoleID := r.Context().Value("userRoleID")

	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, web.ErrIDMustBeenPosInt)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	if tokenUserID != userID && tokenUserRoleID != AdminRoleID {
		web.WriteForbidden(w, nil)
		return
	}

	body := new(model.EditUserBody)
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}
	defer r.Body.Close()

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

// DeleteUser
// @Router /api/v1/users/{id} [delete]
// @Tags Users
// @Param id path int true "User ID"
// @Description Удалить пользователя (Пользователь - только себя, Админ - любого)
// @Security BearerAuth
// @Success 204
// @Failure 403 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
// TODO: удаление должно быть в транзакции с удалением всех связанных элементов.
// Пока остальных сервисов (и таблиц) нет, невозможно доделать.
func (us *UsersPGService) DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer web.PanicRecoverWithSlog(w, us.logger, "users.DeleteUser")

	tokenUserID := r.Context().Value("userID")
	tokenUserRoleID := r.Context().Value("userRoleID")

	userID, err := web.ParseIDFromURL(r, "userID")
	if err != nil {
		if errors.Is(err, web.ErrIDMustBeenPosInt) {
			web.WriteBadRequest(w, err)
			return
		}
		web.WriteServerErrorWithSlog(w, us.logger, err)
		return
	}

	if tokenUserID != userID && tokenUserRoleID != AdminRoleID {
		web.WriteForbidden(w, nil)
		return
	}

	err = us.store.DeleteUser(userID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
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
