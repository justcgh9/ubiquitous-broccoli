package app

import (
	"context"
	"log/slog"

	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
	grpcapp "github.com/justcgh9/discord-clone-friends/internal/app/grpc/server"
	"github.com/justcgh9/discord-clone-friends/internal/config"
	friends "github.com/justcgh9/discord-clone-friends/internal/service"
	"github.com/justcgh9/discord-clone-friends/internal/storage/graph"
	"github.com/justcgh9/discord-clone-friends/internal/storage/postgres"
	storage "github.com/justcgh9/discord-clone-friends/internal/storage/sync"
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
	postgresPath string,
	graphCfg config.GraphStorage,
	loginCfg config.GRPCUserServiceClient,
) *App {

	postgresRepo := postgres.MustConnect(context.Background(), postgresPath)
	graphRepo := graph.NewGraphRepository(graph.MustConnect(
		graphCfg.URI,
		graphCfg.Username,
		graphCfg.Password,
		graphCfg.Realm,
	))

	db := storage.NewRepository(postgresRepo, graphRepo, log)

	
	friendsService := friends.NewService(
		db,
		db,
		db,
	)
	
	loginRequestQueue, err := client.NewLoginByTokenPool(
		log,
		loginCfg.URI,
		loginCfg.NumWorkers,
		loginCfg.QueueSize,
		loginCfg.Timeout,
	)

	if err != nil {
		panic(err)
	}


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