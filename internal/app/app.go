package app

import (
	grpcapp "atlogex/gofoyer/internal/app/grpc"
	"log/slog"
)

type App struct {
	//log        *slog.Logger
	GRPCServer *grpcapp.App
	//port       int
	//tokenTTL   string
}

func New(
	log *slog.Logger,
	//gRPCServer *grpc.Server,
	port int,
	grpcPath string,
	tokenTTL string,
) *App {

	//GRPCServer := grpc.NewServer()
	//auth.Register(gRPCServer)

	grpcApp := grpcapp.New(log, port)

	return &App{
		//log:        log,
		GRPCServer: grpcApp,
		//port:       port,
		//tokenTTL:   tokenTTL,
	}
}
