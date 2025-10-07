package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/goggle-source/grpc-servic/sso/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	// TODO: инициализировать бд

	// TODO: инициализировать сервисный слой

	grpcApp := grpcapp.NewApp(log, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}

}
