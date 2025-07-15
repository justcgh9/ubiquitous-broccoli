package app

import (
	"context"
	"log/slog"
	"time"

	kafka "github.com/justcgh9/discord-clone-kafka"
	grpcapp "github.com/justcgh9/discord-clone-users/internal/app/grpc"
	"github.com/justcgh9/discord-clone-users/internal/service/auth"
	"github.com/justcgh9/discord-clone-users/internal/storage/postgres"
)

type App struct {
	GRPCApp	*grpcapp.App
}

func New(
	log *slog.Logger,
	gRPCPort int,
	storagePath string,
	producer *kafka.Producer,
	tokenTTL time.Duration,
) *App {

	storage := postgres.MustConnect(context.Background(), storagePath)

	authService := auth.New(
		log,
		storage,
		storage,
		storage,
		tokenTTL,
		producer,
	)

	gRPCApp := grpcapp.New(log, authService, gRPCPort)

	return &App{
		GRPCApp: gRPCApp,
	}
}