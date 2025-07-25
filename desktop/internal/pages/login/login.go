package login

import (
	"context"
	"fmt"
	"log/slog"
	"net/mail"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/justcgh9/discord-clone/desktop/internal/appcontext"
	"github.com/justcgh9/discord-clone/desktop/internal/components/background"
	"github.com/justcgh9/discord-clone/desktop/internal/models/user"
	"github.com/justcgh9/discord-clone/desktop/internal/pages/register"
)

func ShowLoginPage(ctx *appcontext.Context, log *slog.Logger) {
	var win fyne.Window
	if ctx.App.ActiveWindow() == nil {
		win = ctx.App.NewWindow("ForkCord")
	} else {
		win = ctx.App.ActiveWindow()
	}

	log = log.With(slog.String("page", "login"))
	ctx.App.SetActiveWindow(win)

	closeWindowChan := make(chan struct{})

	// Form fields
	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	submit := func() {
		emailText := strings.TrimSpace(email.Text)
		passwordText := strings.TrimSpace(password.Text)

		if emailText == "" || passwordText == "" {
			dialog.ShowError(fmt.Errorf("email and password are required"), win)
			return
		}

		if _, err := mail.ParseAddress(emailText); err != nil {
			dialog.ShowError(fmt.Errorf("invalid email address"), win)
			return
		}

		log.Info(fmt.Sprintf("login clicked: %s", emailText))

		usr, token, err := ctx.RPC.Login(context.TODO(), user.NewLoginDTO(
			emailText,
			passwordText,
		))
		if err != nil {
			dialog.ShowError(fmt.Errorf("login failed: %v", err), win)
			return
		}

		ctx.User = usr
		ctx.User.Token = token

		closeWindowChan <- struct{}{}
	}

	// Handle Enter key
	email.OnSubmitted = func(_ string) { submit() }
	password.OnSubmitted = func(_ string) { submit() }

	// Form layout
	loginButton := widget.NewButtonWithIcon("Log In", theme.ConfirmIcon(), submit)
	forgot := widget.NewHyperlink("Forgot your password?", nil)

	registerLink := widget.NewHyperlink("Register", nil)
	registerLink.OnTapped = func() {
		register.ShowRegisterPage(ctx, log, func() {
			ShowLoginPage(ctx, log)
		})
	}

	form := container.NewVBox(
		widget.NewLabelWithStyle("Welcome back!", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("Email *"),
		email,
		widget.NewLabel("Password *"),
		password,
		forgot,
		loginButton,
		widget.NewLabel("Don't have an account?"),
		registerLink,
	)

	// Background & layout
	bg := background.NewBackgroundImage("./media/background-login.png")
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

	// Close window after successful login
	go func() {
		<-closeWindowChan
		fyne.Do(win.Close)
	}()

	win.SetContent(content)
	win.Resize(fyne.NewSize(900, 600))
	win.Show()
}
