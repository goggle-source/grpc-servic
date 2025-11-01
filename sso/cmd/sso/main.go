package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/goggle-source/grpc-servic/sso/internal/app"
	"github.com/goggle-source/grpc-servic/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	log := SetupLogger(cfg.Env)
	log.Info("start server", slog.Int("port", cfg.GRPC.Port))

	application := app.NewApp(log, cfg.GRPC.Port, *cfg, cfg.TokenTTL)

	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("applciation stop")
}

func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
