package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/tokens"
	"github.com/B-Dmitriy/expenses/pgk/web"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
)

const (
	CookieName = "refresh_token"
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
		if errors.Is(err, pgx.ErrNoRows) {
			web.WriteNotFound(w, fmt.Errorf("user with email %s not found", body.Email))
			return
		}
		web.WriteServerErrorWithSlog(w, as.logger, err)
		return
	}

	if !as.passManager.CheckPasswordHash(body.Password, candidate.Password) {
		web.WriteBadRequest(w, fmt.Errorf("login credential wrong"))
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
		if errors.As(err, &storage.ErrNotFound) {
			web.WriteNotFound(w, nil)
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

func (as *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	web.PanicRecoverWithSlog(w, as.logger, "auth.Refresh")

	c, err := r.Cookie(CookieName)
	if err != nil {
		web.WriteBadRequest(w, fmt.Errorf("refresh token in cookie not found"))
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
