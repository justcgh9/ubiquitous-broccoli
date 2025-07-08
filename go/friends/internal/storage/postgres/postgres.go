package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justcgh9/discord-clone-friends/internal/models"
	"github.com/justcgh9/discord-clone-friends/internal/storage"
)

type PgxInterface interface {
	BeginTx(context.Context, pgx.TxOptions) (pgx.Tx, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type friendRepo struct {
	db PgxInterface
}

func MustConnect(ctx context.Context, connStr string) *friendRepo {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatalf("%s: %v", op, err)
	}

	return &friendRepo{
		db: conn,
	}
}

func (r *friendRepo) Close() {
	r.db.Close()
}

func (r *friendRepo) SendRequest(ctx context.Context, fromUserID, toUserID string) error {
	const op = "storage.postgres.SendRequest"

	_, err := r.db.Exec(ctx, `
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES ($1, $2, 'PENDING')
		ON CONFLICT DO NOTHING
	`, fromUserID, toUserID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *friendRepo) AcceptRequest(ctx context.Context, userID, targetID string) error {
	const op = "storage.postgres.AcceptRequest"

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	cmdTag, err := tx.Exec(ctx, `
		UPDATE friendships SET status = 'ACCEPTED'
		WHERE user_id = $1 AND friend_id = $2 AND status = 'PENDING'
	`, targetID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, storage.ErrNoPendingRequest)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES ($1, $2, 'ACCEPTED')
		ON CONFLICT (user_id, friend_id) DO UPDATE SET status = 'ACCEPTED'
	`, userID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return tx.Commit(ctx)
}

func (r *friendRepo) DenyRequest(ctx context.Context, userID, targetID string) error {
	const op = "storage.postgres.DenyRequest"

	_, err := r.db.Exec(ctx, `
		DELETE FROM friendships
		WHERE user_id = $1 AND friend_id = $2 AND status = 'PENDING'
	`, targetID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *friendRepo) RemoveFriend(ctx context.Context, userID, targetID string) error {
	const op = "storage.postgres.RemoveFriend"

	_, err := r.db.Exec(ctx, `
		DELETE FROM friendships
		WHERE ((user_id = $1 AND friend_id = $2)
		   OR (user_id = $2 AND friend_id = $1))
		   AND status = 'ACCEPTED'
	`, userID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *friendRepo) BlockUser(ctx context.Context, userID, targetID string) error {
	const op = "storage.postgres.BlockUser"

	_, err := r.db.Exec(ctx, `
		INSERT INTO friendships (user_id, friend_id, status)
		VALUES ($1, $2, 'BLOCKED')
		ON CONFLICT (user_id, friend_id) DO UPDATE SET status = 'BLOCKED'
	`, userID, targetID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *friendRepo) ListFriends(ctx context.Context, userID string) ([]models.Friendship, error) {
	const op = "storage.postgres.ListFriends"

	rows, err := r.db.Query(ctx, `
		SELECT friend_id, status
		FROM friendships
		WHERE user_id = $1 AND status = 'ACCEPTED'
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var friends []models.Friendship
	for rows.Next() {
		var friendID, status string
		if err := rows.Scan(&friendID, &status); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		friends = append(friends, models.Friendship{
			UserID:   userID,
			FriendID: friendID,
			Status:   models.FriendStatus(status),
		})
	}
	return friends, nil
}
