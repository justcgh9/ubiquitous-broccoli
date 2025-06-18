package app

import (
	"log/slog"

	"fyne.io/fyne/v2"
	fyneapp "github.com/justcgh9/discord-clone/desktop/internal/app/fyne"
	grpcapp "github.com/justcgh9/discord-clone/desktop/internal/app/grpc"
	"github.com/justcgh9/discord-clone/desktop/internal/appcontext"
)

type App struct {
	app *fyneapp.App
	ctx *appcontext.Context
}

func Run (
	log *slog.Logger,
	app fyne.App,
	srvAddr string,
) *App {

	fyneApp := fyneapp.New(
		log,
		app,
	)

	grpc, err := grpcapp.Connect(
		log,
		srvAddr,
	)

	if err != nil {
		msg := "failed to connect to gRPC server"
		log.Error(msg, slog.String("err", err.Error()))
		fyneApp.FailRun(msg)
		return nil
	}

	ctx := appcontext.Context{
		App: app,
		Log: log,
		RPC: grpc,
	}

	fyneApp.Run(ctx)

	return &App{
		app: fyneApp,
		ctx: &ctx,
	}
}