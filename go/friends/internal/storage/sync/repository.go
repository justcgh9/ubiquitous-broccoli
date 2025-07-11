package storage

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/justcgh9/discord-clone-friends/internal/models"
)

type Graph interface {
	CreateUser(ctx context.Context, id, handle string) error
	SendFriendRequest(ctx context.Context, fromID, toID string) error
	AcceptFriendRequest(ctx context.Context, userID, requesterID string) error
	RemoveFriend(ctx context.Context, userID, friendID string) error
	BlockUser(ctx context.Context, userID, targetID string) error
	ListFriends(ctx context.Context, userID string) ([]string, error)
	ListMutualFriends(ctx context.Context, userA, userB string) ([]string, error)
}

type Postgres interface {
	SendRequest(ctx context.Context, fromUserID, toUserID string) error
	AcceptRequest(ctx context.Context, userID, targetID string) error
	DenyRequest(ctx context.Context, userID, targetID string) error
	RemoveFriend(ctx context.Context, userID, targetID string) error
	BlockUser(ctx context.Context, userID, targetID string) error
	ListFriends(ctx context.Context, userID string) ([]models.Friendship, error)
}

type syncOpType string

const (
	opSendRequest   syncOpType = "SendRequest"
	opAcceptRequest syncOpType = "AcceptRequest"
	opRemoveFriend  syncOpType = "RemoveFriend"
	opBlockUser     syncOpType = "BlockUser"
)

type syncJob struct {
	op      syncOpType
	from    string
	to      string
	attempt int
}

type Repository struct {
	pg        Postgres
	graph     Graph
	jobs      chan syncJob
	wg        sync.WaitGroup
	retries   int
	log       *slog.Logger
	component string
}

func NewRepository(pg Postgres, graph Graph, logger *slog.Logger) *Repository {
	r := &Repository{
		pg:        pg,
		graph:     graph,
		jobs:      make(chan syncJob, 1000),
		retries:   3,
		log:       logger.With("layer", "storage", "component", "sync-repo"),
		component: "sync-repo",
	}

	r.log.Info("starting sync repository background worker")
	r.wg.Add(1)
	go r.worker()

	return r
}

func (r *Repository) Close() {
	r.log.Info("shutting down sync repository")
	close(r.jobs)
	r.wg.Wait()
}

func (r *Repository) worker() {
	r.log.Info("sync worker started")

	defer r.wg.Done()
	for job := range r.jobs {
		r.log.Info("processing sync job", "op", job.op, "from", job.from, "to", job.to)

		for job.attempt = 1; job.attempt <= r.retries; job.attempt++ {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			err := r.perform(ctx, job)
			cancel()

			if err == nil {
				r.log.Info("sync job succeeded", "op", job.op, "from", job.from, "to", job.to, "attempt", job.attempt)
				break
			}

			r.log.Error("sync job failed", "op", job.op, "from", job.from, "to", job.to, "attempt", job.attempt, "err", err)
			time.Sleep(time.Duration(job.attempt) * time.Second)
		}
	}
}

func (r *Repository) perform(ctx context.Context, job syncJob) error {
	switch job.op {
	case opSendRequest:
		return r.graph.SendFriendRequest(ctx, job.from, job.to)
	case opAcceptRequest:
		return r.graph.AcceptFriendRequest(ctx, job.from, job.to)
	case opRemoveFriend:
		return r.graph.RemoveFriend(ctx, job.from, job.to)
	case opBlockUser:
		return r.graph.BlockUser(ctx, job.from, job.to)
	default:
		return errors.New("unknown sync job type")
	}
}

// --- Public API ---

func (r *Repository) SendRequest(ctx context.Context, from, to string) error {
	const op = "storage.sync.SendRequest"

	r.log.Info("writing to postgres", "op", op, "from", from, "to", to)
	if err := r.pg.SendRequest(ctx, from, to); err != nil {
		r.log.Error("postgres write failed", "op", op, "err", err)
		return fmt.Errorf("%s (postgres): %w", op, err)
	}

	job := syncJob{op: opSendRequest, from: from, to: to}
	r.jobs <- job
	r.log.Info("enqueued sync job", "op", job.op, "from", job.from, "to", job.to)

	return nil
}

func (r *Repository) AcceptRequest(ctx context.Context, userID, requesterID string) error {
	const op = "storage.sync.AcceptRequest"

	r.log.Info("writing to postgres", "op", op, "userID", userID, "requesterID", requesterID)
	if err := r.pg.AcceptRequest(ctx, userID, requesterID); err != nil {
		r.log.Error("postgres write failed", "op", op, "err", err)
		return fmt.Errorf("%s (postgres): %w", op, err)
	}

	job := syncJob{op: opAcceptRequest, from: userID, to: requesterID}
	r.jobs <- job
	r.log.Info("enqueued sync job", "op", job.op, "from", job.from, "to", job.to)

	return nil
}

func (r *Repository) RemoveFriend(ctx context.Context, userID, friendID string) error {
	const op = "storage.sync.RemoveFriend"

	r.log.Info("writing to postgres", "op", op, "userID", userID, "friendID", friendID)
	if err := r.pg.RemoveFriend(ctx, userID, friendID); err != nil {
		r.log.Error("postgres write failed", "op", op, "err", err)
		return fmt.Errorf("%s (postgres): %w", op, err)
	}

	job := syncJob{op: opRemoveFriend, from: userID, to: friendID}
	r.jobs <- job
	r.log.Info("enqueued sync job", "op", job.op, "from", job.from, "to", job.to)

	return nil
}

func (r *Repository) BlockUser(ctx context.Context, userID, targetID string) error {
	const op = "storage.sync.BlockUser"

	r.log.Info("writing to postgres", "op", op, "userID", userID, "targetID", targetID)
	if err := r.pg.BlockUser(ctx, userID, targetID); err != nil {
		r.log.Error("postgres write failed", "op", op, "err", err)
		return fmt.Errorf("%s (postgres): %w", op, err)
	}

	job := syncJob{op: opBlockUser, from: userID, to: targetID}
	r.jobs <- job
	r.log.Info("enqueued sync job", "op", job.op, "from", job.from, "to", job.to)

	return nil
}

func (r *Repository) DenyRequest(ctx context.Context, userID, requesterID string) error {
	const op = "storage.sync.DenyRequest"
	r.log.Info("writing to postgres", "op", op, "userID", userID, "requesterID", requesterID)
	return r.pg.DenyRequest(ctx, userID, requesterID)
}

func (r *Repository) ListFriends(ctx context.Context, userID string) ([]models.Friendship, error) {
	const op = "storage.sync.ListFriends"
	r.log.Info("reading from postgres", "op", op, "userID", userID)
	return r.pg.ListFriends(ctx, userID)
}
