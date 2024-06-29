package postgres

import (
	"context"
	"fmt"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/jackc/pgx/v5"
)

type PGStorage struct {
	Conn  *pgx.Conn
	Utils *PGUtils
}

func NewStorage(settings config.StorageSettings) (*PGStorage, error) {
	// ?sslmode=disable
	urlExample := fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		settings.DBDriver,
		settings.DBUser,
		settings.DBPass,
		settings.DBHost,
		settings.DBPort,
		settings.DBName,
	)

	conn, err := pgx.Connect(context.Background(), urlExample)
	if err != nil {
		return nil, err
	}

	return &PGStorage{
		Conn:  conn,
		Utils: &PGUtils{},
	}, nil
}

func (s *PGStorage) CloseConnection(ctx context.Context) error {
	return s.Conn.Close(ctx)
}
