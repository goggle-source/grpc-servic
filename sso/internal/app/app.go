package app

import (
	"log/slog"
	"time"

	grpcapp "github.com/goggle-source/grpc-servic/sso/internal/app/grpc"
	"github.com/goggle-source/grpc-servic/sso/internal/config"
	"github.com/goggle-source/grpc-servic/sso/internal/services/auth"
	"github.com/goggle-source/grpc-servic/sso/internal/storage/postgresql"
)

type App struct {
	GRPCServer *grpcapp.App
}

func NewApp(log *slog.Logger, grpcPort int, cfg config.Config, tokenTTL time.Duration) *App {

	db, err := postgresql.New(cfg)
	if err != nil {
		panic(err)
	}

	auth := auth.New(log, db, db, db, tokenTTL)

	grpcApp := grpcapp.NewApp(log, grpcPort, auth)

	return &App{
		GRPCServer: grpcApp,
	}

}
