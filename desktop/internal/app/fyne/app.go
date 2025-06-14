package fyneapp

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"github.com/justcgh9/discord-clone/desktop/internal/pages/login"
)

type App struct {
	app fyne.App
	log *slog.Logger
}

func New(
	log *slog.Logger,
	app fyne.App,
) *App {
	return &App{
		log: slog.With(
			slog.String("client", "fyne app"),
		),
		app: app,
	}
}

func (a *App) Run() {
	
	login.ShowLoginPage(
		a.app,
		a.log,
	)

	a.app.Run()
}