package users

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
)

type UsersStorage struct {
	db *postgres.PGStorage
}

func NewUsersStorage(db *postgres.PGStorage) storage.UsersStore {
	return &UsersStorage{
		db: db,
	}
}

func (s *UsersStorage) GetUsersList(page, limit int, search string) ([]*model.UserInfo, error) {
	users := make([]*model.UserInfo, 0)

	offset, err := s.db.Utils.GetOffset(page, limit)
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Conn.Query(context.Background(), `
		SELECT * FROM users 
			WHERE LOWER(email) LIKE CONCAT('%', $1::text,'%') 
			   OR LOWER(login) LIKE CONCAT('%', $1::text,'%')
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, search, limit, offset)
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

	err := s.db.Conn.QueryRow(context.Background(), "SELECT * FROM users WHERE id = $1;", id).Scan(
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
		if errors.As(err, &pgx.ErrNoRows) {
			return nil, storage.ErrNotFound
		}
		return nil, err
	}

	return convertUserToUserInfo(user), nil
}

func (s *UsersStorage) GetUserByEmail(email string) (*model.User, error) {
	user := new(model.User)

	err := s.db.Conn.QueryRow(context.Background(), "SELECT * FROM users WHERE email = $1;", email).Scan(
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

	return user, nil
}

func (s *UsersStorage) CreateUser(body *model.CreateUserBody) error {
	_, err := s.db.Conn.Exec(context.Background(), `
		INSERT INTO users (login, email, password) 
		VALUES ($1, $2, $3);`,
		body.Login,
		body.Email,
		body.Password,
	)
	if err != nil {
		return err
	}

	return nil
}

func (s *UsersStorage) EditUser(id int, user *model.EditUserBody) error {
	res, err := s.db.Conn.Exec(context.Background(), "UPDATE users SET login=$1,email=$2 WHERE id = $3;",
		user.Login,
		user.Email,
		id,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *UsersStorage) DeleteUser(id int) error {
	res, err := s.db.Conn.Exec(context.Background(), "DELETE FROM users WHERE id = $1;", id)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
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
