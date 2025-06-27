package fyneapp

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"github.com/justcgh9/discord-clone/desktop/internal/appcontext"
	errorpage "github.com/justcgh9/discord-clone/desktop/internal/pages/error"
	"github.com/justcgh9/discord-clone/desktop/internal/pages/login"
)

type App struct {
	fyne.App
	activeWindow fyne.Window
	log 		 *slog.Logger
}

func New(
	log *slog.Logger,
	app fyne.App,
) *App {
	return &App{
		log: slog.With(
			slog.String("client", "fyne app"),
		),
		App: app,
	}
}

func (a *App) Start(
	ctx *appcontext.Context,
) {

	login.ShowLoginPage(
		ctx,
		a.log,
	)
	a.Run()
}

func (a *App) RenderError(message string) {
	a.log.Error("failure", slog.String("err", message))
	errorpage.ShowErrorWindow(a, message)
}

func (a *App) FailRun(message string) {
	a.RenderError(message)
	a.Run()
}

func (a *App) ActiveWindow() fyne.Window {
	return  a.activeWindow
}

func (a *App) SetActiveWindow(w fyne.Window) {
	a.activeWindow = w
}