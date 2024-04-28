package mysql

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	ID              int
	Username        string
	Password        string
	Email           string
	Email_confirmed bool
}

func (s *Storage) GetUser(id int) (*User, error) {
	const op = "storage.mysql.GetUser"

	rows, err := s.DB.Query("SELECT * FROM users WHERE id = ? LIMIT 1;", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFoundError
		}
	}

	user := new(User)

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Email_confirmed); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return user, nil
}

func (s *Storage) GetUsers(limit, offset int) ([]User, error) {
	const op = "storage.mysql.GetUsers"

	rows, err := s.DB.Query("SELECT * FROM users LIMIT ? OFFSET ?;", limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	users := make([]User, 0, 0)

	for rows.Next() {
		user := new(User)
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Email_confirmed); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, *user)
	}

	return users, err
}
