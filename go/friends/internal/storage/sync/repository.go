package storage

import (
    "context"
    "errors"
    "fmt"
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
    pg      Postgres
    graph   Graph
    jobs    chan syncJob
    wg      sync.WaitGroup
    retries int
}

func NewRepository(pg Postgres, graph Graph) *Repository {
    r := &Repository{
        pg:      pg,
        graph:   graph,
        jobs:    make(chan syncJob, 1000),
        retries: 3,
    }

    r.wg.Add(1)
    go r.worker()

    return r
}

func (r *Repository) Close() {
    close(r.jobs)
    r.wg.Wait()
}

func (r *Repository) worker() {
    defer r.wg.Done()
    for job := range r.jobs {
        for job.attempt = 1; job.attempt <= r.retries; job.attempt++ {
            ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
            err := r.perform(ctx, job)
            cancel()
            if err == nil {
                break
            }
            time.Sleep(time.Duration(job.attempt) * time.Second) // backoff
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
        return errors.New("unknown job type")
    }
}

// Public API

func (r *Repository) SendRequest(ctx context.Context, from, to string) error {
    const op = "storage.sync.SendRequest"

    if err := r.pg.SendRequest(ctx, from, to); err != nil {
        return fmt.Errorf("%s (postgres): %w", op, err)
    }

    r.jobs <- syncJob{op: opSendRequest, from: from, to: to}
    return nil
}

func (r *Repository) AcceptRequest(ctx context.Context, userID, requesterID string) error {
    const op = "storage.sync.AcceptRequest"

    if err := r.pg.AcceptRequest(ctx, userID, requesterID); err != nil {
        return fmt.Errorf("%s (postgres): %w", op, err)
    }

    r.jobs <- syncJob{op: opAcceptRequest, from: userID, to: requesterID}
    return nil
}

func (r *Repository) RemoveFriend(ctx context.Context, userID, friendID string) error {
    const op = "storage.sync.RemoveFriend"

    if err := r.pg.RemoveFriend(ctx, userID, friendID); err != nil {
        return fmt.Errorf("%s (postgres): %w", op, err)
    }

    r.jobs <- syncJob{op: opRemoveFriend, from: userID, to: friendID}
    return nil
}

func (r *Repository) BlockUser(ctx context.Context, userID, targetID string) error {
    const op = "storage.sync.BlockUser"

    if err := r.pg.BlockUser(ctx, userID, targetID); err != nil {
        return fmt.Errorf("%s (postgres): %w", op, err)
    }

    r.jobs <- syncJob{op: opBlockUser, from: userID, to: targetID}
    return nil
}

func (r *Repository) DenyRequest(ctx context.Context, userID, requesterID string) error {
    return r.pg.DenyRequest(ctx, userID, requesterID)
}

func (r *Repository) ListFriends(ctx context.Context, userID string) ([]models.Friendship, error) {
    return r.pg.ListFriends(ctx, userID)
}
