package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
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

	loginRequestQueue, err := client.NewLoginByTokenPool(
		log,
		cfg.UsersClient.URI,
		cfg.UsersClient.NumWorkers,
		cfg.UsersClient.QueueSize,
		cfg.UsersClient.Timeout,
	)

	if err != nil {
		panic(err)
	}

	resC := loginRequestQueue.Enqueue(context.TODO(), "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBfaWQiOjEsImVtYWlsIjoianVzdGNvb2xlc3RnaXJhZmZlOUBnbWFpbC5jb20iLCJleHAiOjE3NTQ0ODA5NjIsInVpZCI6M30.8mOyl5UlNQ7au7cMTOmt4xHIkBuUGCxCDUSY1uYX484")

	res := <- resC

	fmt.Println(res)

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