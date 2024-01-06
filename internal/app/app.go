package app

import (
	grpcapp "auth-sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseUri string,
	tokenTTL time.Duration,
) *App {
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
