package users

import (
	"context"
	"errors"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/jackc/pgx/v5"
)

const (
	DefaultUserRoleID = 10
)

var (
	ErrNotFound = errors.New("not found")
)

type UsersStorage struct {
	db *pgx.Conn
}

func NewUsersStorage(db *pgx.Conn) *UsersStorage {
	return &UsersStorage{
		db: db,
	}
}

func (s *UsersStorage) GetUsersList() ([]*model.UserInfo, error) {
	users := make([]*model.UserInfo, 0)
	rows, err := s.db.Query(context.Background(), "SELECT * FROM users;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		user := model.User{}

		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.EmailConfirmed,
			&user.Password,
			&user.RoleID,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, convertUserToUserInfo(&user))
	}

	return users, nil
}

func (s *UsersStorage) GetUser(id int) (*model.UserInfo, error) {
	user := new(model.User)

	stmt, err := s.db.Prepare(context.Background(), "getUser", "SELECT * FROM users WHERE id = $1;")
	if err != nil {
		return nil, err
	}

	err = s.db.QueryRow(context.Background(), stmt.SQL, id).Scan(
		&user.ID,
		&user.Login,
		&user.Email,
		&user.EmailConfirmed,
		&user.Password,
		&user.RoleID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return convertUserToUserInfo(user), nil
}

func (s *UsersStorage) CreateUser(body *model.CreateUserBody) error {
	_, err := s.db.Exec(context.Background(), "INSERT INTO users (login, email, password, role_id) VALUES ($1, $2, $3, $4);",
		body.Login,
		body.Email,
		body.Password,
		DefaultUserRoleID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStorage) EditUser(id int, user *model.EditUserBody) error {
	res, err := s.db.Exec(context.Background(), "UPDATE users SET login=$1,email=$2 WHERE id = $3;",
		user.Login,
		user.Email,
		id,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (s *UsersStorage) DeleteUser(id int) error {
	res, err := s.db.Exec(context.Background(), "DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func convertUserToUserInfo(user *model.User) *model.UserInfo {
	return &model.UserInfo{
		ID:             user.ID,
		Login:          user.Login,
		Email:          user.Email,
		EmailConfirmed: user.EmailConfirmed,
		RoleID:         user.RoleID,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
