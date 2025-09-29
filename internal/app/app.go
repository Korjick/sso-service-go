package app

import (
	grpcapp "github.com/Korjick/sso-service-go/internal/app/grpc"
	"github.com/Korjick/sso-service-go/internal/services/auth"
	"github.com/Korjick/sso-service-go/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	logger *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(logger, storage, storage, storage, tokenTTL)

	return &App{
		GRPCServer: grpcapp.New(logger, authService, grpcPort),
	}
}
