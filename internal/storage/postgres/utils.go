package postgres

import (
	"errors"

	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/jackc/pgx/v5/pgconn"
)

func NewPGUtils() *PGUtils {
	return &PGUtils{}
}

type PGUtils struct{}

func (pgu *PGUtils) GetOffset(page, limit int) (int, error) {
	switch {
	case page < 1:
		return 0, storage.ErrPageMustBeenGreaterThanOne
	case page == 1:
		return 0, nil
	case limit < 1:
		return 0, storage.ErrLimitMustBeenGreaterThanOne
	default:
		return (page - 1) * limit, nil
	}
}

func (pgu *PGUtils) CheckConstrainError(e error) (bool, error) {
	var pgErr *pgconn.PgError
	if errors.As(e, &pgErr) {
		switch pgErr.ConstraintName {
		case "users_unique_login":
			return true, storage.ErrUsersUniqueLogin
		case "users_unique_email":
			return true, storage.ErrUsersUniqueEmail
		case "users_empty_login":
			return true, storage.ErrUsersEmptyLogin
		case "users_empty_email":
			return true, storage.ErrUsersEmptyEmail
		case "users_empty_password":
			return true, storage.ErrUsersEmptyPassword
		case "categories_user_category_name":
			return true, storage.ErrCategoryUniqueName
		default:
			return false, e
		}
	}
	return false, e
}
