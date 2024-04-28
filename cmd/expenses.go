package main

import (
	"log/slog"
	"os"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/logger"
	"github.com/B-Dmitriy/expenses/internal/storage/mysql"

	muxHTTP "github.com/B-Dmitriy/expenses/internal/server/http"
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

	httpServer := muxHTTP.NewServer(config.HTTPServer, storage, logger)
	err = httpServer.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
