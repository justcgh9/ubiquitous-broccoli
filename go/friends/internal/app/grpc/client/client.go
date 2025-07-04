package client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/justcgh9/discord-clone-proto/gen/go/users"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LoginByTokenPool struct {
	requests chan loginRequest
	log      *slog.Logger
}

type loginRequest struct {
	token   string
	resultC chan loginResult
}

type loginResult struct {
	user *users.LoginByTokenResponse
	err  error
}

func NewLoginByTokenPool(
	log *slog.Logger,
	target string,
	numWorkers int,
	queueSize int,
) (*LoginByTokenPool, error) {
	pool := &LoginByTokenPool{
		requests: make(chan loginRequest, queueSize),
		log:      log,
	}

	for i := 0; i < numWorkers; i++ {
		conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		))
		if err != nil {
			return nil, fmt.Errorf("worker %d: failed to dial: %w", i, err)
		}
		client := users.NewUserServiceClient(conn)

		go pool.worker(i, client)
	}

	return pool, nil
}

func (p *LoginByTokenPool) worker(id int, client users.UserServiceClient) {
	p.log.Info("worker started", slog.Int("id", id))
	for req := range p.requests {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		resp, err := client.LoginByToken(ctx, &users.LoginByTokenRequest{
			AccessToken: req.token,
		})

		req.resultC <- loginResult{
			user: resp,
			err:  err,
		}
	}
}
