package http

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/B-Dmitriy/expenses/internal/config"
)

type HTTPServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(cfg config.HTTPSettings, logger *slog.Logger) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := http.NewServeMux()

	handler := initRoutes(r)

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

func (s *HTTPServer) Run() {
	if err := s.server.ListenAndServe(); err != nil {
		s.logger.Error(fmt.Sprintf("start server error: %v", err))
		os.Exit(1)
		return
	}
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
