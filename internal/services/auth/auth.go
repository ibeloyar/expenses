package auth

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/tokens"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/go-playground/validator/v10"
)

const (
	CookieName = "refresh_token"
)

var (
	ErrCredentialWrong      = errors.New("login credential wrong")
	ErrRefreshTokenNotFound = errors.New("refresh token in cookie not found")
)

type AuthService struct {
	logger        *slog.Logger
	usersStorage  storage.UsersStore
	validator     *validator.Validate
	tokensStorage storage.TokensStore
	utils         storage.ServiceUtils
	tokensManager *tokens.TokensManager
	passManager   *password.PasswordManager
}

func NewAuthService(
	logger *slog.Logger,
	utils storage.ServiceUtils,
	validator *validator.Validate,
	usersStorage storage.UsersStore,
	tokensStorage storage.TokensStore,
	tokensManager *tokens.TokensManager,
	passManager *password.PasswordManager,
) *AuthService {
	return &AuthService{
		utils:         utils,
		logger:        logger,
		validator:     validator,
		passManager:   passManager,
		usersStorage:  usersStorage,
		tokensStorage: tokensStorage,
		tokensManager: tokensManager,
	}
}

// Login
// @Router /api/v1/login [post]
// @Tags Authentication
// @Description Вход под учётной записью
// @Param request body model.LoginCredentials false "query params"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (as *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	web.PanicRecoverWithSlog(w, as.logger, "auth.Login")

	body := new(model.LoginCredentials)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}
	defer r.Body.Close()

	err = as.validator.Struct(body)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		web.WriteBadRequest(w, errs)
		return
	}

	candidate, err := as.usersStorage.GetUserByEmail(body.Email)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	if !as.passManager.CheckPasswordHash(body.Password, candidate.Password) {
		web.WriteBadRequest(w, ErrCredentialWrong)
		return
	}

	tkns, err := as.tokensManager.GenerateTokens(candidate.ID, candidate.RoleID)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	isTokenExist, err := as.tokensStorage.CheckToken(candidate.ID)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	if isTokenExist {
		err := as.tokensStorage.ChangeToken(candidate.ID, tkns.RefreshToken)
		if err != nil {
			web.WriteServerErrorWithSlog(w, as.logger, err)
			return
		}

	} else {
		err := as.tokensStorage.CreateToken(candidate.ID, tkns.RefreshToken)
		if err != nil {
			web.WriteServerErrorWithSlog(w, as.logger, err)
			return
		}
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    tkns.RefreshToken,
		HttpOnly: true,
		MaxAge:   2 * 24 * 3600, // 14 days
	})

	web.WriteOK(w, &model.LoginResponse{
		AccessToken:  tkns.AcceptToken,
		RefreshToken: tkns.RefreshToken,
		Login:        candidate.Login,
		UserID:       candidate.ID,
		UserRoleID:   candidate.RoleID,
	})
}

// Registration
// @Router /api/v1/registration [post]
// @Tags Authentication
// @Description Регистрация пользователя
// @Param request body model.RegistrationData false "query params"
// @Success 201
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (as *AuthService) Registration(w http.ResponseWriter, r *http.Request) {
	web.PanicRecoverWithSlog(w, as.logger, "auth.Registration")

	body := new(model.RegistrationData)

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}
	defer r.Body.Close()

	hashPass, err := as.passManager.HashPassword(body.Password)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	err = as.usersStorage.CreateUser(&model.CreateUserBody{
		Email:    body.Email,
		Login:    body.Login,
		Password: hashPass,
	})

	if err != nil {
		if isConstrain, e := as.utils.CheckConstrainError(err); isConstrain {
			web.WriteBadRequest(w, e)
			return
		}
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	web.WriteCreated(w, nil)
}

// Logout
// @Router /api/v1/logout [post]
// @Tags Authentication
// @Description Выход из учётной записи
// @Security BearerAuth
// @Success 204
// @Failure 401 {object} web.WebError
// @Failure 404 {object} web.WebError
// @Failure 500 {object} web.WebError
func (as *AuthService) Logout(w http.ResponseWriter, r *http.Request) {
	web.PanicRecoverWithSlog(w, as.logger, "auth.Logout")

	bearer, err := as.passManager.GetAuthorizationHeader(r)
	if err != nil {
		web.WriteUnauthorized(w, err)
		return
	}

	userInfo, err := as.tokensManager.VerifyJWTToken(bearer)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	err = as.tokensStorage.DeleteToken(userInfo.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			web.WriteNotFound(w, storage.ErrNotFound)
			return
		}
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		HttpOnly: true,
		MaxAge:   -1,
	})

	web.WriteNoContent(w, nil)
}

// Refresh
// @Router /api/v1/refresh [post]
// @Tags Authentication
// @Description Перегенерировать токен пользователя
// @Security BearerAuth
// @Success 200 {object} tokens.Tokens
// @Failure 400 {object} web.WebError
// @Failure 500 {object} web.WebError
func (as *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	web.PanicRecoverWithSlog(w, as.logger, "auth.Refresh")

	c, err := r.Cookie(CookieName)
	if err != nil {
		web.WriteBadRequest(w, ErrRefreshTokenNotFound)
		return
	}

	userInfo, err := as.tokensManager.VerifyJWTToken(c.Value)
	if err != nil {
		web.WriteBadRequest(w, err)
		return
	}

	tkns, err := as.tokensManager.GenerateTokens(userInfo.UserID, userInfo.UserRoleID)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	err = as.tokensStorage.ChangeToken(userInfo.UserID, tkns.RefreshToken)
	if err != nil {
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    tkns.RefreshToken,
		HttpOnly: true,
		MaxAge:   2 * 24 * 3600, // 14 days
	})

	web.WriteOK(w, &tokens.Tokens{
		AcceptToken:  tkns.AcceptToken,
		RefreshToken: tkns.RefreshToken,
	})
}
