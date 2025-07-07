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
	ctx context.Context
	token   string
	resultC chan loginResult
}

type loginResult struct {
	User *users.LoginByTokenResponse
	Err  error
}

func NewLoginByTokenPool(
	log *slog.Logger,
	target string,
	numWorkers int,
	queueSize int,
	timeout time.Duration,
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

		ctx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()
	
		_, err = client.Ping(ctx, &users.PingRequest{})
		if err != nil {
			log.Error("could not connect", slog.Int("worker id", i), slog.String("err", err.Error()))
			return nil, fmt.Errorf("worker %d: failed to dial: %w", i, err)
		}
	

		go pool.worker(i, client, timeout)
	}

	return pool, nil
}

func (p *LoginByTokenPool) worker(
	id int,
	client users.UserServiceClient,
	timeout time.Duration,
) {
	p.log.Info("worker started", slog.Int("id", id))
	for req := range p.requests {
		select {
		case <-req.ctx.Done():
			req.resultC <- loginResult{User: nil, Err: req.ctx.Err()}
			continue
		default:
		}

		ctx, cancel := context.WithTimeout(req.ctx, timeout)
		resp, err := client.LoginByToken(ctx, &users.LoginByTokenRequest{
			AccessToken: req.token,
		})
		cancel()

		req.resultC <- loginResult{
			User: resp,
			Err:  err,
		}
	}
}

func (p *LoginByTokenPool) Enqueue(ctx context.Context, token string) <-chan loginResult {
	resultC := make(chan loginResult, 1)

	select {
	case p.requests <- loginRequest{ctx: ctx, token: token, resultC: resultC}:
	case <-ctx.Done():
		resultC <- loginResult{User: nil, Err: ctx.Err()}
	}

	return resultC
}