package app

import (
	"log/slog"

	grpcapp "github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/app/grpc"
	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/service"
	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/storage/sqlite"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	srvc := service.New(log, storage, storage)

	grpcApp := grpcapp.New(srvc, log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
