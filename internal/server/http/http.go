package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/B-Dmitriy/expenses/internal/config"
	"github.com/B-Dmitriy/expenses/internal/storage/postgres"
	"github.com/B-Dmitriy/expenses/pgk/password"
	"github.com/B-Dmitriy/expenses/pgk/tokens"

	_ "github.com/B-Dmitriy/expenses/api"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HTTPServer struct {
	server *http.Server
	logger *slog.Logger
}

func NewServer(
	env string,
	cfg config.HTTPSettings,
	logger *slog.Logger,
	db *postgres.PGStorage,
	tm *tokens.TokensManager,
	pm *password.PasswordManager,
) *HTTPServer {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	r := http.NewServeMux()

	handler := initRoutes(r, logger, db, tm, pm)

	if env != config.ProductionENV {
		handler.HandleFunc("GET /swagger/*", httpSwagger.Handler())
	}

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
		// From go doc https://pkg.go.dev/net/http#Server.Shutdown
		//
		// When Shutdown is called, Serve, ListenAndServe, and ListenAndServeTLS
		// immediately return ErrServerClosed. Make sure the program doesn't exit
		// and waits instead for Shutdown to return.
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	}

	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
