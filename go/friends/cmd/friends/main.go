package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/justcgh9/discord-clone-friends/internal/config"
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

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<- done
	log.Info("user service stopped")
}

func setupLogger(env string) *slog.Logger {
	var l *slog.Logger 

	switch env {
	case envLocal:
		l = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
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