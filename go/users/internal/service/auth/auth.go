package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/justcgh9/discord-clone-users/internal/lib/jwt"
	"github.com/justcgh9/discord-clone-users/internal/models"
	"github.com/justcgh9/discord-clone-users/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, handle string, passHash []byte) (int64, error)
}

type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

func New(
	log *slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	email string,
	password string,
	appID int,
) (string, models.UserDTO, error) {
	const op = "Auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	log.Info("login attempt")

	user, err := a.usrProvider.User(ctx, email)
	if errors.Is(err, storage.ErrUserNotFound) {
		log.Warn("user not found", slog.String("err", err.Error()))
		return "", models.UserDTO{}, fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	if err != nil {
		log.Error("failed to get user", slog.String("err", err.Error()))
		return "", models.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		log.Warn("invalid credentials", slog.String("err", err.Error()))
		return "", models.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		log.Error("failed to get app", slog.String("err", err.Error()))
		return "", models.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(
		user,
		app,
		a.tokenTTL,
	)

	if err != nil {
		log.Error("failed to generate token", slog.String("err", err.Error()))
		return "", models.UserDTO{}, fmt.Errorf("%s: %w", op, err)
	}

	return token, models.NewDTOFromUser(user), nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, handle string, pass string) (int64, error) {
	const op = "Auth.RegisterNewUser"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
		slog.String("pass", pass),
	)

	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", slog.String("err", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.SaveUser(ctx, email, handle, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("err", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "Auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
