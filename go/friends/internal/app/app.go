package app

import (
	"context"
	"log/slog"

	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
	grpcapp "github.com/justcgh9/discord-clone-friends/internal/app/grpc/server"
	friends "github.com/justcgh9/discord-clone-friends/internal/service"
)


type AuthValidator interface {
	Verify(ctx context.Context, token string) <-chan client.LoginResult
}

type App struct {
	GRPCApp *grpcapp.App
	Auth AuthValidator
}

func New(
	log *slog.Logger,
	gRPCPort int,
	friendsService *friends.Service,
	loginRequestQueue AuthValidator,
) *App {

	return &App{
		GRPCApp: grpcapp.New(
			log,
			friendsService,
			loginRequestQueue,
			gRPCPort,
		),
		Auth: loginRequestQueue,
	}
}