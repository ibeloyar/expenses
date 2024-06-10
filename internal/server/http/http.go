package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/jackc/pgx/v5"
)

type HTTPServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(cfg config.HTTPSettings, logger *slog.Logger, db *pgx.Conn) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := http.NewServeMux()

	handler := initRoutes(r, logger, db)

	return &HTTPServer{
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			WriteTimeout: time.Duration(cfg.Timeout) * time.Second,
			ReadTimeout:  time.Duration(cfg.IddleTimout) * time.Second,
		},
		logger: logger,
	}
}

func (s *HTTPServer) Run() error {
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(fmt.Sprintf("start server error: %v", err))
		return err
	}
	return nil
}

func (s *HTTPServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
