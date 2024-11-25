package main

import (
	"atlogex/gofoyer/internal/app"
	"atlogex/gofoyer/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envProduction  = "production"
	envDevelopment = "development"
	envLocal       = "local"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
	logger := setupLogger(cfg.Env)
	logger.Error("Start", slog.String("env", cfg.Env))

	application := app.New(logger, cfg.GRPCPort, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	signalOs := <-stop

	application.GRPCServer.Stop()

	logger.Info(
		"Application Stopped for ",
		slog.String("env", cfg.Env),
		slog.String("signal", signalOs.String()))
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envProduction:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	case envDevelopment:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	}

	return logger
}
