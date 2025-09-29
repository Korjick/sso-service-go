package grpcapp

import (
	"fmt"
	authgrpc "github.com/Korjick/sso-service-go/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	logger     *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

// New creates a new gRPC server app.
func New(logger *slog.Logger, authService authgrpc.Auth, port int) *App {
	gRPCServer := grpc.NewServer()
	authgrpc.Register(gRPCServer, authService)

	return &App{
		logger:     logger,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "grpcapp.Run"

	logger := a.logger.With(slog.String("op", op))

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: failed to listen: %w", op, err)
	}

	logger.Info("gRPC server is running",
		slog.String("addr", l.Addr().String()),
		slog.Int("port", a.port),
	)

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: failed to serve: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.logger.With(op, slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
