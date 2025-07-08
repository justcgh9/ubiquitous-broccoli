package app

import (
	"context"

	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
)


type AuthValidator interface {
	Verify(ctx context.Context, token string) <-chan client.LoginResult
}

type FriendsApp interface {

}

type App struct {
	Srv FriendsApp
	Auth AuthValidator
}