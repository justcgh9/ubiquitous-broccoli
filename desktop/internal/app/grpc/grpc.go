package grpcapp

import (
	"context"
	"log/slog"
	"time"

	"github.com/justcgh9/discord-clone-proto/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client users.UserServiceClient
}

func Connect(
	log *slog.Logger,
) (*GRPCClient, error) {

	log = log.With(
		slog.String("client", "gRPC"),
	)

	log.Info("gRPC connection attempt")
	conn, err := grpc.NewClient("localhost:44044", grpc.WithTransportCredentials(
		insecure.NewCredentials(),
	))
	if err != nil {
		log.Error("could not connect", slog.String("err", err.Error()))
		return nil, err
	}

	client := users.NewUserServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	_, err = client.Ping(ctx, &users.PingRequest{})
	if err != nil {
		log.Error("could not connect", slog.String("err", err.Error()))
		return nil, err
	}

	log.Info("gRPC connection succeeded")

	return &GRPCClient{
		conn:   conn,
		client: client,
	}, nil
}
