package main

import (
	"auth-sso/internal/app"
	"auth-sso/internal/config"
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "production"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info("Starting application", slog.String("env", cfg.Env))

	application := app.New(
		log,
		cfg.GRPC.Port,
		cfg.Database.Uri,
		cfg.Database.DatabaseName,
		cfg.TokenTTL,
	)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	receivedSignal := <-stop
	log.Info("Stopping application", slog.String("signal", receivedSignal.String()))

	if err := application.Storage.Close(context.Background()); err != nil {
		log.Error("Failed to close the MongoDB connection", slog.String("error", err.Error()))
	}

	log.Info("MongoDB connection is closed.")

	application.GRPCServer.Stop()

	log.Info("Application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
