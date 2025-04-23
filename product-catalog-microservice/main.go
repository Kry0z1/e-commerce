package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Kry0z1/e-commerce/logger/handlers/slogpretty"
	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/app"
	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/config"
)

var (
	localStr = "local"
	prodStr  = "prod"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	logger.Info("Server gracefully died")
}

func setupLogger(level string) *slog.Logger {
	switch level {
	case localStr:
		return slog.New(slogpretty.NewPrettyHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case prodStr:
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		return slog.Default()
	}
}
