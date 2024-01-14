package main

import (
	"auth-sso/internal/app"
	"auth-sso/internal/config"
	"auth-sso/internal/tasks"
	"context"
	"github.com/hibiken/asynq"
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
		cfg.Redis.Address,
	)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	log.Info("Starting Asynq client", slog.String("redisAddr", cfg.Redis.Address))
	redisClient := asynq.RedisClientOpt{Addr: cfg.Redis.Address}
	asynqServer := asynq.NewServer(redisClient, asynq.Config{Concurrency: 10})

	mux := asynq.NewServeMux()
	tasks.SetupTaskHandlers(mux)

	go func() {
		if err := asynqServer.Run(mux); err != nil {
			log.Error("Asynq server stopped with error", slog.String("error", err.Error()))
		}
	}()

	receivedSignal := <-stop
	log.Info("Stopping application", slog.String("signal", receivedSignal.String()))

	asynqServer.Shutdown()
	log.Info("Asynq server shut down")

	if err := application.Storage.Close(context.Background()); err != nil {
		log.Error("Failed to close the MongoDB connection", slog.String("error", err.Error()))
	}

	log.Info("MongoDB connection is closed.")

	application.GRPCServer.Stop()

	if err := application.AsynqClient.Close(); err != nil {
		log.Error("Failed to close the Asynq client", slog.String("error", err.Error()))
	}

	log.Info("Asynq client is closed.")

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
