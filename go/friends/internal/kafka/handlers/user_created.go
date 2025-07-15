package handlers

import (
	"context"
	"fmt"
	"log/slog"

	kafka "github.com/justcgh9/discord-clone-kafka"
)

func UserCreated(
	createUser func(context.Context, string, string) error,
) func(context.Context, []byte) error {
	return  func(ctx context.Context, data []byte) error {
		event, err := kafka.ParseUserCreatedEvent(data)
		if err != nil {
			return fmt.Errorf("error parsing event: %s", err.Error())
		}

		slog.Info("read a user created event", slog.Any("event", event))

		return createUser(ctx, event.UserID, event.Handle)
	}
}