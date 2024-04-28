package mysql

import (
	"database/sql"
	"fmt"

	"github.com/B-Dmitriy/expenses/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	DB *sql.DB
}

func NewMySQLStorage(cfg *config.StorageSettings) (*Storage, error) {
	const op = "storage.mysql"

	db, err := sql.Open(cfg.DBDriver, cfg.DBUser+":"+cfg.DBPass+"@/"+cfg.DBName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{DB: db}, nil
}
