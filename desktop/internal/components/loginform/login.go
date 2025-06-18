package loginform

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
	"github.com/justcgh9/discord-clone/desktop/internal/models/user"
)

func NewLoginForm(ctx appcontext.Context, log *slog.Logger, done chan<- struct{}) *fyne.Container {
	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	formWindow := ctx.App.NewWindow("Login")

	submit := func() {
		// Validate inputs
		emailText := strings.TrimSpace(email.Text)
		passwordText := strings.TrimSpace(password.Text)

		if emailText == "" || passwordText == "" {
			dialog.ShowError(fmt.Errorf("email and password are required"), formWindow)
			return
		}

		if _, err := mail.ParseAddress(emailText); err != nil {
			dialog.ShowError(fmt.Errorf("invalid email address"), formWindow)
			return
		}

		log.Info(fmt.Sprintf("login clicked: %s %s", emailText, passwordText))

		usr, token, err := ctx.RPC.Login(context.TODO(), user.NewLoginDTO(
			emailText,
			passwordText,
		))
		if err != nil {
			dialog.ShowError(fmt.Errorf("login failed: %v", err), formWindow)
			return
		}

		_, _ = usr, token

		done <- struct{}{}
	}

	loginButton := widget.NewButtonWithIcon("Log In", theme.ConfirmIcon(), submit)

	// Handle Enter key for both inputs
	email.OnSubmitted = func(_ string) { submit() }
	password.OnSubmitted = func(_ string) { submit() }

	forgot := widget.NewHyperlink("Forgot your password?", nil)
	register := widget.NewHyperlink("Register", nil)

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
		register,
	)

	formWindow.SetContent(form)

	return form
}
