package tokens

import (
	"context"
	"errors"
	"github.com/B-Dmitriy/expenses/internal/model"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/jackc/pgx/v5"
)

type TokensStorage struct {
	db *postgres.PGStorage
}

func NewTokensStorage(db *postgres.PGStorage) storage.TokensStore {
	return &TokensStorage{
		db: db,
	}
}

func (ts *TokensStorage) GetTokenByUserID(userID int) (*model.Token, error) {
	token := new(model.Token)
	err := ts.db.Conn.QueryRow(
		context.Background(),
		"SELECT * FROM refresh_tokens WHERE user_id = $1",
		userID,
	).Scan(
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (ts *TokensStorage) CheckToken(userID int) (bool, error) {
	token := new(model.Token)

	err := ts.db.Conn.QueryRow(
		context.Background(),
		"SELECT * FROM refresh_tokens WHERE user_id = $1;",
		userID,
	).Scan(
		&token.UserID,
		&token.Token,
		&token.CreatedAt,
	)
	if err != nil {
		if errors.As(err, &pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, err
}

func (ts *TokensStorage) CreateToken(userID int, token string) error {
	_, err := ts.db.Conn.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (user_id, token) VALUES ($1, $2)",
		userID, token,
	)
	if err != nil {
		return err
	}

	return nil
}

func (ts *TokensStorage) ChangeToken(userID int, token string) error {
	res, err := ts.db.Conn.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET token = $1 WHERE user_id = $2;",
		token, userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (ts *TokensStorage) DeleteToken(userID int) error {
	res, err := ts.db.Conn.Exec(
		context.Background(),
		"DELETE FROM refresh_tokens WHERE user_id = $1;",
		userID,
	)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return storage.ErrNotFound
	}

	return nil
}
