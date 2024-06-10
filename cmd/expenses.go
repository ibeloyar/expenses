package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/storage"
	"github.com/B-Dmitriy/expenses/pgk/logger"

	server "github.com/B-Dmitriy/expenses/internal/server/http"
)

func main() {
	cfg := config.MustLoad()

	lgr := logger.NewLogger(cfg.ENV)
	lgr.Info("logger initialized", slog.String("env", cfg.ENV), slog.String("port", strconv.Itoa(cfg.HTTPServer.Port)))

	s, err := storage.NewStorage(cfg.Storage)
	if err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}

	if err := server.NewServer(cfg.HTTPServer, lgr, s.Conn).Run(); err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}
}
