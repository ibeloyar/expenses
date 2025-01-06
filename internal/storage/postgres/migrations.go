package postgres

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/ibeloyar/expenses/internal/config"

	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5"
)

func (s *PGStorage) RunMigration(settings config.StorageSettings) error {
	dbString := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable",
		settings.DBDriver,
		settings.DBUser,
		settings.DBPass,
		settings.DBHost,
		settings.DBPort,
		settings.DBName,
	)

	db, err := sql.Open("postgres", dbString)
	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("migrations", "postgres", driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil {
		return err
	}

	return nil
}
