package app

import (
	grpcapp "auth-sso/internal/app/grpc"
	"auth-sso/internal/services/auth"
	"auth-sso/internal/services/identity"
	"auth-sso/internal/storage/mongodb"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	Storage    *mongodb.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	databaseUri string,
	database string,
	tokenTTL time.Duration,
) *App {
	client, err := mongodb.New(databaseUri, database)
	if err != nil {
		panic(err)
	}

	log.Info("MongoDB connection is successful.")

	authService := auth.New(log, client, client, client, tokenTTL)
	identityService := identity.New(log)
	grpcApp := grpcapp.New(log, authService, identityService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		Storage:    client,
	}
}
