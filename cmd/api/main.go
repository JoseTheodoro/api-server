package main

import (
	"apiserver/internal/config"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"apiserver/internal/config/observability/logging"
	hhttp "apiserver/internal/http"
	"apiserver/internal/http/handlers"
	"apiserver/internal/repository/postgres"
	"apiserver/internal/routes"
	"apiserver/internal/services"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

	ctx := context.Background()
	cfg := config.NewConfig()

	shut, err := cfg.StartTrace(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = shut(context.Background())
	}()

	l := logging.NewLogger()
	slog.SetDefault(l.Log)

	conn := postgres.NewConnection(os.Getenv("DATABASE_DSN"))
	db, err := conn.Connect(ctx)
	if err != nil {
		slog.Error("unable connect db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	userRepositoryPostgres := postgres.NewUserRepositoryPostgres(db)
	userService := services.NewUserService(userRepositoryPostgres)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()
	routes.Register(mux, routes.Handlers{User: userHandler})

	instrumented := otelhttp.NewHandler(mux, "http.server")

	s := hhttp.NewServer(":9090", instrumented)
	errCh := make(chan error, 1)
	slog.Info("server started successful", "port", ":9090")
	go func() {
		if err := s.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	ctxNotify, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case err := <-errCh:
		slog.Error("unknown error received", "err", err)
	case <-ctxNotify.Done():
		slog.Info("shutdown signal recivied")

	}

	ctxsht, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.Shutdown(ctxsht); err != nil {
		slog.Error("error on graceful shutdown", "err", err)
	}

}
