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

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func initTrace(ctx context.Context) (func(context.Context) error, error) {
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpoint("alloy:4318"), otlptracehttp.WithInsecure())
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName("api-go")))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil

}

func main() {

	ctx := context.Background()

	shut, err := initTrace(ctx)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = shut(context.Background())
	}()

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
