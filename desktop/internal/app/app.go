package app

import (
	"log/slog"

	"fyne.io/fyne/v2"
	fyneapp "github.com/justcgh9/discord-clone/desktop/internal/app/fyne"
	grpcapp "github.com/justcgh9/discord-clone/desktop/internal/app/grpc"
)

type App struct {
	log *slog.Logger
	app *fyneapp.App
	rpc *grpcapp.GRPCClient
}

func New(
	log *slog.Logger,
	app fyne.App,
) *App {
	return &App{
		log: log,
		app: fyneapp.New(
			log,
			app,
		),
	}
}

func (a *App) Run() {
	grpc, err := grpcapp.Connect(a.log)
	if err != nil {
		msg := "failed to connect to gRPC server"
		a.log.Error(msg, slog.String("err", err.Error()))
		a.app.FailRun(msg)
	} else {
		a.rpc = grpc
		a.app.Run()
	}
}
