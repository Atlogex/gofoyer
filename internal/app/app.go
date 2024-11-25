package app

import (
	"atlogex/gofoyer/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
)

type App struct {
	log        *slog.Logger
	GRPCServer *grpc.Server
	port       int
	tokenTTL   string
}

func New(log *slog.Logger, gRPCServer *grpc.Server, port int, tokenTTL string) *App {

	GRPCServer := grpc.NewServer()
	auth.Register(gRPCServer)

	return &App{
		log:        log,
		GRPCServer: GRPCServer,
		port:       port,
		tokenTTL:   tokenTTL,
	}
}
