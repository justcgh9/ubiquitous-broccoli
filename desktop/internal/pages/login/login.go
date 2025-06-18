package login

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"github.com/justcgh9/discord-clone/desktop/internal/appcontext"
	"github.com/justcgh9/discord-clone/desktop/internal/components/background"
	"github.com/justcgh9/discord-clone/desktop/internal/components/loginform"
)

func ShowLoginPage(ctx appcontext.Context, log *slog.Logger) {
	win := ctx.App.NewWindow("ForkCord")
	log = log.With(
		slog.String("page", "login"),
	)

	closeWindowChan := make(chan struct{})
	
	bg := background.NewBackgroundImage("./media/background-login.png")
	form := loginform.NewLoginForm(
		ctx,
		log,
		closeWindowChan,
	)

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

	
	go func() {
		<- closeWindowChan
		fyne.Do(win.Close)
	} ()

	win.SetContent(content)
	win.Resize(fyne.NewSize(900, 600))
	win.Show()
}