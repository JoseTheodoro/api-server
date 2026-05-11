package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	hhttp "apiserver/internal/http"
	"apiserver/internal/http/handlers"
	"apiserver/internal/repository/memory"
	"apiserver/internal/repository/postgres"
	"apiserver/internal/routes"
	"apiserver/internal/services"
)

func main() {

	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(log)

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

	productRepository := memory.NewProductRepositoryMemory()
	productService := services.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	mux := http.NewServeMux()
	routes.Register(mux, routes.Handlers{User: userHandler, Product: productHandler})

	s := hhttp.NewServer(":9090", mux)
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
