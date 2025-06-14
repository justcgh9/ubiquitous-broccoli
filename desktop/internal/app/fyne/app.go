package fyneapp

import (
	"log/slog"

	"fyne.io/fyne/v2"
	errorpage "github.com/justcgh9/discord-clone/desktop/internal/pages/error"
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

func (a *App) RenderError(message string) {
	a.log.Error("failure", slog.String("err", message))
	errorpage.ShowErrorWindow(a.app, message)
}

func (a *App) FailRun(message string) {
	a.RenderError(message)
	a.app.Run()
}
