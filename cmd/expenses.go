package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/logger"
	"github.com/B-Dmitriy/expenses/internal/storage/mysql"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config := config.MustLoad()

	logger := logger.SetupLogger(config.ENV)
	logger.Info("logger initialized", slog.String("env", config.ENV))

	storage, err := mysql.NewMySQLStorage(&config.Storage)
	if err != nil {
		logger.Error("database connect error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	users, err := storage.GetUsers(5, 2)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, v := range users {
		fmt.Printf("User: %v\n", v)
	}

	user, err := storage.GetUser(2)
	if err != nil {
		logger.Error(err.Error())
	}
	fmt.Printf("User: %v\n", user)
	// TODO: implement server
	// TODO: Run app
}
