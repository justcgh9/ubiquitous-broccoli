package authgrpc

import (
	"context"
	"errors"

	"github.com/justcgh9/discord-clone-proto/gen/go/users"
	"github.com/justcgh9/discord-clone-users/internal/service/auth"
	"github.com/justcgh9/discord-clone-users/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (string, error)
	RegisterNewUser(ctx context.Context,
		email string,
		handle string,
		pass string,
	) (int64, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	users.UnimplementedUserServiceServer
	auth Auth
}

func RegisterServer(server *grpc.Server, auth Auth) {
	users.RegisterUserServiceServer(server, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *users.LoginRequest,
) (*users.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if errors.Is(err, auth.ErrInvalidCredentials) {
		return nil, status.Error(codes.InvalidArgument, "invalid email or password")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to log in")
	}

	return &users.LoginResponse{
		AccessToken: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *users.RegisterRequest,
) (*users.RegisterResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "username is required")
	}

	id, err := s.auth.RegisterNewUser(
		ctx,
		req.GetEmail(),
		req.GetUsername(),
		req.GetPassword(),
	)

	if errors.Is(err, storage.ErrUserExists) {
		return nil, status.Error(codes.AlreadyExists, "user already exists")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &users.RegisterResponse{
		UserId: id,
	}, nil
}

func (s *serverAPI) Ping(
	ctx context.Context,
	req *users.PingRequest,
) (*users.PongResponse, error) {
	return &users.PongResponse{}, nil
}