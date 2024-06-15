package postgres

import (
	"context"
	"fmt"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/jackc/pgx/v5"
	//"github.com/golang-migrate/migrate/v5"
	//"github.com/golang-migrate/migrate/v5/database/postgres"
	//"github.com/golang-migrate/migrate/v5/source/iofs"
)

////go:embed migrations/*.sql
//var expensesdb embed.FS
//
//func MigrateSchema(db *sql.DB, _ *config.StorageSettings) error {
//	source, err := iofs.New(expensesdb, "preferences")
//	if err != nil {
//		return err
//	}
//	conf := new(config.StorageSettings)
//	conf.DBName = "postgres"
//	target, err := postgres.WithInstance(db, conf)
//	if err != nil {
//		return err
//	}
//	m, err := migrate.NewWithInstance("iofs", source, conf.DBName, target)
//	if err != nil {
//		return err
//	}
//	err = m.Up()
//	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
//		return err
//	}
//	return source.Close()
//}

type PGStorage struct {
	Conn  *pgx.Conn
	Utils *PGUtils
}

func NewStorage(settings config.StorageSettings) (*PGStorage, error) {
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
