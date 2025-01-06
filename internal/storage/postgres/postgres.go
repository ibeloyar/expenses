package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/ibeloyar/expenses/internal/config"
	"github.com/jackc/pgx/v5"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)

type PGStorage struct {
	Conn       *pgx.Conn
	Utils      *PGUtils
	connString string
}

func NewStorage(settings config.StorageSettings) (*PGStorage, error) {
	connString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		settings.DBDriver,
		settings.DBUser,
		settings.DBPass,
		settings.DBHost,
		settings.DBPort,
		settings.DBName,
	)

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil, err
	}

	return &PGStorage{
		Conn:       conn,
		Utils:      &PGUtils{},
		connString: connString,
	}, nil
}

func (s *PGStorage) CloseConnection(ctx context.Context) error {
	return s.Conn.Close(ctx)
}

//go:embed migrations/*.sql
var fs embed.FS

func (s *PGStorage) MigrateSchema() error {
	data, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", data, s.connString)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		// Error db no change for init migration
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}

	return data.Close()
}
