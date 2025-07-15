package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/justcgh9/discord-clone-friends/internal/app"
	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
	"github.com/justcgh9/discord-clone-friends/internal/config"
	"github.com/justcgh9/discord-clone-friends/internal/kafka/handlers"
	friends "github.com/justcgh9/discord-clone-friends/internal/service"
	"github.com/justcgh9/discord-clone-friends/internal/storage/graph"
	"github.com/justcgh9/discord-clone-friends/internal/storage/postgres"
	storage "github.com/justcgh9/discord-clone-friends/internal/storage/sync"
	kafka "github.com/justcgh9/discord-clone-kafka"
)

const (
	envLocal = "local"
	envTest = "test"
	envProd = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	postgresRepo := postgres.MustConnect(context.Background(), cfg.StoragePath)
	defer postgresRepo.Close()

	graphRepo := graph.NewGraphRepository(graph.MustConnect(
		cfg.GraphStorage.URI,
		cfg.GraphStorage.Username,
		cfg.GraphStorage.Password,
		cfg.GraphStorage.Realm,
	))
	defer graphRepo.Close(context.Background())

	db := storage.NewRepository(postgresRepo, graphRepo, log)
	defer db.Close()

	
	friendsService := friends.NewService(
		db,
		db,
		db,
	)
	
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

	defer loginRequestQueue.Close()


	log.Info("starting up the user service")

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	app := app.New(
		log,
		cfg.GRPCSrv.Port,
		friendsService,
		loginRequestQueue,
	)

	go app.GRPCApp.MustRun()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		kafka.TopicUserCreated,
		cfg.Kafka.GroupID,
		handlers.UserCreated(graphRepo.CreateUser),
		cfg.Kafka.MinBytes,
		cfg.Kafka.MaxBytes,
		cfg.Kafka.MaxWait,
	)

	go consumer.Start(ctx)
	defer consumer.Close()

	<- done
	app.GRPCApp.Stop()
	defer log.Info("user service stopped")
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