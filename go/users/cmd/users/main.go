package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/justcgh9/discord-clone-users/internal/app"
	"github.com/justcgh9/discord-clone-users/internal/config"
	"github.com/justcgh9/discord-clone-users/internal/lib/logger/handlers/pretty"
)

const (
	envLocal = "local"
	envTest = "test"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting up the user service")

	myApp := app.New(
		log,
		cfg.GRPC.Port,
		cfg.StoragePath,
		cfg.TokenTTL,
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		myApp.GRPCApp.MustRun()
	} ()

	<- done
	myApp.GRPCApp.Stop()

	log.Info("user service stopped")
}

func setupLogger(env string) *slog.Logger {
	var l *slog.Logger 

	switch env {
	case envLocal:
		l = slog.New(
			pretty.NewPrettyHandler(os.Stdout, pretty.PrettyHandlerOptions{
				ShowCaller: true,
				IndentJSON: true,
			}),
		)
	case envTest:
		l = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		l = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return l
}