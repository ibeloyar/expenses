package storage

import "errors"

var (
	ErrNotFound                    = errors.New("not found")
	ErrPageMustBeenGreaterThanOne  = errors.New("page must been greater than one")
	ErrLimitMustBeenGreaterThanOne = errors.New("limit must been greater than one")
	ErrUsersUniqueLogin            = errors.New("user must have unique login")
	ErrUsersUniqueEmail            = errors.New("user must have unique email")
	ErrUsersEmptyLogin             = errors.New("login is required")
	ErrUsersEmptyEmail             = errors.New("email is required")
	ErrUsersEmptyPassword          = errors.New("password is required")
)
