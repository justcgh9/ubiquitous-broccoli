package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/justcgh9/discord-clone-users/internal/models"
	"github.com/justcgh9/discord-clone-users/internal/storage"
)

const (
	uniqueConstraintViolationCode = "23505"
)

type PgxInterface interface {
	Begin(context.Context) (pgx.Tx, error)
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Close()
}

type Storage struct {
	conn PgxInterface
}

func New(ctx context.Context, connStr string) *Storage {
	const op = "storage.postgres.New"

	conn, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatalf("%s %v", op, err)
	}

	return &Storage{
		conn: conn,
	}
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, handle string, passHash []byte) (int64, error) {
	const op = "storage.postgres.SaveUser"

	query := `INSERT INTO users(email, handle, pass_hash) VALUES($1, $2, $3) RETURNING id`
	var id int64
	err := s.conn.QueryRow(ctx, query, email, handle, passHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueConstraintViolationCode { 
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.postgres.User"

	query := `SELECT id, email, handle, pass_hash FROM users WHERE email = $1`
	var user models.User
	err := s.conn.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.Handle, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.postgres.App"

	query := `SELECT id, name, secret FROM apps WHERE id = $1`
	var app models.App
	err := s.conn.QueryRow(ctx, query, id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	return app, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	query := `SELECT is_admin FROM users WHERE id = $1`
	var isAdmin bool
	err := s.conn.QueryRow(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return isAdmin, nil
}
