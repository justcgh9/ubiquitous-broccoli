package appcontext

import (
	"log/slog"

	"fyne.io/fyne/v2"
	grpcapp "github.com/justcgh9/discord-clone/desktop/internal/app/grpc"
)

type Context struct {
	App fyne.App
	Log *slog.Logger
	RPC *grpcapp.GRPCClient
}
