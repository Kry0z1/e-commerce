package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/Kry0z1/e-commerce/sso-microservice/internal/app/grpc"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/services/auth"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(authService, log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
