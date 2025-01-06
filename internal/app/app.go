package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ibeloyar/expenses/internal/config"
	"github.com/ibeloyar/expenses/internal/storage/postgres"
	"github.com/ibeloyar/expenses/pgk/logger"
	"github.com/ibeloyar/expenses/pgk/password"
	"github.com/ibeloyar/expenses/pgk/tokens"

	server "github.com/ibeloyar/expenses/internal/server/http"
)

func Run(cfg *config.Config) {

	lgr := logger.NewLogger(cfg.ENV)
	lgr.Info("logger initialized", slog.String("env", cfg.ENV), slog.String("port", strconv.Itoa(cfg.HTTPServer.Port)))

	pm := password.New(cfg.Security.PassCost)
	lgr.Info("password manager initialized")

	tm := tokens.New(cfg.Security.JWTSecret)
	lgr.Info("tokens manager initialized")

	// TODO: Зачем здесь возвращать err?
	// Без Storage приложение не работоспособно
	store, err := postgres.NewStorage(cfg.Storage)
	if err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}

	// TODO: Зачем здесь возвращать err?
	// Без миграций приложение не работоспособно
	err = store.MigrateSchema()
	if err != nil {
		lgr.Error(err.Error())
		os.Exit(1)
	}

	srv := server.NewServer(cfg, lgr, store, tm, pm)

	go func() {
		err = srv.Run()
		if err != nil {
			lgr.Error(fmt.Sprintf("start server error: %v", err))
			os.Exit(1)
		}
	}()
	lgr.Info("server started", slog.Int("port", cfg.HTTPServer.Port))

	shutdownTimeout := time.Second * time.Duration(cfg.HTTPServer.ShutdownTimeout)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case <-ctx.Done():
		// Shutdown the server gracefully
		fmt.Println("\nShutting down HTTP server gracefully...")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancelShutdown()

		err := srv.Stop(shutdownCtx)
		if err != nil {
			fmt.Printf("HTTP server shutdown error: %s\n", err)
		}
		fmt.Println("- http server is stopped")

		err = store.CloseConnection(shutdownCtx)
		if err != nil {
			fmt.Printf("database server shutdown error: %s\n", err)
		}
		fmt.Println("- database connection is closed")
	}
}
