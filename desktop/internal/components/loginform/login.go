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
	"github.com/justcgh9/discord-clone/desktop/internal/pages/register"
)

func NewLoginForm(ctx *appcontext.Context, 
	log *slog.Logger,
	done chan<- struct{},
	back func(),
	) *fyne.Container {
	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	submit := func() {
		// Validate inputs
		emailText := strings.TrimSpace(email.Text)
		passwordText := strings.TrimSpace(password.Text)

		if emailText == "" || passwordText == "" {
			dialog.ShowError(fmt.Errorf("email and password are required"), ctx.App.ActiveWindow())
			return
		}

		if _, err := mail.ParseAddress(emailText); err != nil {
			dialog.ShowError(fmt.Errorf("invalid email address"), ctx.App.ActiveWindow())
			return
		}

		log.Info(fmt.Sprintf("login clicked: %s %s", emailText, passwordText))

		usr, token, err := ctx.RPC.Login(context.TODO(), user.NewLoginDTO(
			emailText,
			passwordText,
		))
		if err != nil {
			dialog.ShowError(fmt.Errorf("login failed: %v", err), ctx.App.ActiveWindow())
			return
		}

		ctx.User = usr
		ctx.User.Token = token

		done <- struct{}{}
	}

	loginButton := widget.NewButtonWithIcon("Log In", theme.ConfirmIcon(), submit)

	// Handle Enter key for both inputs
	email.OnSubmitted = func(_ string) { submit() }
	password.OnSubmitted = func(_ string) { submit() }

	forgot := widget.NewHyperlink("Forgot your password?", nil)
	reg := widget.NewHyperlink("Register", nil)
	reg.OnTapped = func() {
		register.ShowRegisterPage(ctx, log, back)
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
		reg,
	)

	return form
}
