package appcontext

import (
	"log/slog"

	"fyne.io/fyne/v2"
	grpcapp "github.com/justcgh9/discord-clone/desktop/internal/app/grpc"
	"github.com/justcgh9/discord-clone/desktop/internal/models/user"
)

type App interface {
	fyne.App
	ActiveWindow() fyne.Window
	SetActiveWindow(w fyne.Window)
}

type Context struct {
	App 	App
	Log 	*slog.Logger
	RPC 	*grpcapp.GRPCClient
	User 	user.User
}
