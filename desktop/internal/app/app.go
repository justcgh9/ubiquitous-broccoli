package app

import (
	"log/slog"

	"fyne.io/fyne/v2"
	fyneapp "github.com/justcgh9/discord-clone/desktop/internal/app/fyne"
)

type App struct {
	log *slog.Logger
	app *fyneapp.App
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
	a.app.Run()
}