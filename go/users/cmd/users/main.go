package main

import (
	"log/slog"
	"os"

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