package grpcapp

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/justcgh9/discord-clone-proto/gen/go/users"
	"github.com/justcgh9/discord-clone/desktop/internal/models/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client users.UserServiceClient
	log    *slog.Logger
}

const APP_ID = 1

func Connect(
	log *slog.Logger,
	target string,
) (*GRPCClient, error) {

	log = log.With(
		slog.String("client", "gRPC"),
	)

	log.Info("gRPC connection attempt")
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(
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
		log: log,
	}, nil
}

func (c *GRPCClient) Login (
	ctx context.Context, 
	usr user.LoginDTO,
) (user.User, string, error) {
	const op = "grpc.Login"

	log := c.log.With(
		slog.String("op", op),
		slog.String("email", usr.Email),
	)

	log.Info("login attempt")

	resp, err := c.client.Login(
		ctx,
		&users.LoginRequest{
			Email: usr.Email,
			Password: usr.Password,
			AppId: APP_ID,
		},
	)

	if err != nil {
		//TODO: add correct error handling
		log.Warn("error while logging in", slog.String("err", err.Error()))
		return user.User{}, "", err
	}

	log.Info(
		"login successful",
		slog.String("token", resp.GetAccessToken()),
		slog.Any("user", resp.GetUser()),
	)
	

	return user.New(
		fmt.Sprintf("%d", resp.User.UserId),
		resp.User.GetEmail(),
		resp.User.GetHandle(),
	), resp.GetAccessToken(), nil
}

func (c *GRPCClient) Register (
	ctx context.Context,
	usr user.RegisterDTO,
) (int64, error) {
	const op = "grpc.Register"

	log := c.log.With(
		slog.String("op", op),
		slog.String("email", usr.Email),
		slog.String("handle", usr.Handle),
	)

	log.Info("registration attempt")

	resp, err := c.client.Register(
		ctx,
		&users.RegisterRequest{
			Email: usr.Email,
			Username: usr.Handle,
			Password: usr.Password,
		},
	)

	if err != nil {
		//TODO: add correct error handling
		log.Warn("error during registration", slog.String("err", err.Error()))
		return 0, err
	}

	log.Info(
		"login successful",
		slog.Int64("user_id", resp.GetUserId()),
	)

	return resp.GetUserId(), nil
}