package storage

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type PGUtils struct{}

func (pgu *PGUtils) GetOffset(page, limit int) (int, error) {
	switch {
	case page < 1:
		return 0, ErrPageMustBeenGreaterThanOne
	case page == 1:
		return 0, nil
	case limit < 1:
		return 0, ErrLimitMustBeenGreaterThanOne
	default:
		return (page - 1) * limit, nil
	}
}

func (pgu *PGUtils) CheckPGConstrainError(e error) (bool, error) {
	var pgErr *pgconn.PgError
	if errors.As(e, &pgErr) {
		switch pgErr.ConstraintName {
		case "users_unique_login":
			return true, ErrUsersUniqueLogin
		case "users_unique_email":
			return true, ErrUsersUniqueEmail
		case "users_empty_login":
			return true, ErrUsersEmptyLogin
		case "users_empty_email":
			return true, ErrUsersEmptyEmail
		case "users_empty_password":
			return true, ErrUsersEmptyPassword
		default:
			return false, e
		}
	}
	return false, e
}
