package grpcapp

import (
	authgrpc "atlogex/gofoyer/internal/grpc/auth"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	authService authgrpc.Auth,
	port int,
) *App {

	fmt.Println("authService", authService)
	GRPCServer := grpc.NewServer()
	authgrpc.Register(GRPCServer, authService)
	fmt.Println("GRPCServer", GRPCServer)

	return &App{
		log:        log,
		gRPCServer: GRPCServer,
		port:       port,
	}
}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	const operation = "grpcapp.Run"

	log := app.log.With(slog.String("operation", operation))
	log.Info("Starting gRPC server", slog.Int("port", app.port))

	listener, error := net.Listen("tcp", fmt.Sprintf("localhost:%d", app.port))
	if error != nil {
		log.Error("Failed to listen", operation, error)

		return fmt.Errorf("failed to listen: %w", error)
	}

	if err := app.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("GRPC server started - ", slog.Int("port", app.port), listener.Addr().String())

	return nil
}

func (app *App) Stop() {
	const operation = "grpcapp.Stop"
	app.log.With(slog.String("operation", operation)).
		Info("Stopping gRPC server", slog.Int("port", app.port))

	app.gRPCServer.GracefulStop()
}
