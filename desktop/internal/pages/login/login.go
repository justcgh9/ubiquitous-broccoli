package login

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/justcgh9/discord-clone/desktop/internal/components/background"
	"github.com/justcgh9/discord-clone/desktop/internal/components/loginform"
)

func ShowLoginPage(a fyne.App, log *slog.Logger) {
	win := a.NewWindow("ForkCord")

	bg := background.NewBackgroundImage("./media/background-login.png")
	form := loginform.NewLoginForm(log.With(
		slog.String("page", "login"),
	))

	card := container.NewPadded(form)
	cardBG := background.NewCardBackground(
		theme.Color(theme.ColorNameBackground),
		fyne.NewSize(320, 400),
	)
	loginCard := container.NewStack(cardBG, card)

	content := container.NewStack(
		bg,
		container.NewCenter(loginCard),
	)

	win.SetContent(content)
	win.Resize(fyne.NewSize(900, 600))
	win.Show()
}