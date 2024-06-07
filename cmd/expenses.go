package main

import (
	"log/slog"
	"strconv"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/pgk/logger"

	server "github.com/B-Dmitriy/expenses/internal/server/http"
)

func main() {
	cfg := config.MustLoad()

	lgr := logger.SetupLogger(cfg.ENV)
	lgr.Info("logger initialized", slog.String("env", cfg.ENV), slog.String("port", strconv.Itoa(cfg.HTTPServer.Port)))

	server.NewServer(cfg.HTTPServer, lgr).Run()
}
