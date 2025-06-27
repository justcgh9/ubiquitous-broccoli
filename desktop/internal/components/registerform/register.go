package registerform

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

func NewRegisterForm(ctx *appcontext.Context, log *slog.Logger, done chan<- struct{}, back func()) *fyne.Container {
	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	username := widget.NewEntry()
	username.SetPlaceHolder("Username")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	submit := func() {
		emailText := strings.TrimSpace(email.Text)
		usernameText := strings.TrimSpace(username.Text)
		passwordText := strings.TrimSpace(password.Text)

		if emailText == "" || usernameText == "" || passwordText == "" {
			dialog.ShowError(fmt.Errorf("all fields are required"), ctx.App.ActiveWindow())
			return
		}

		if _, err := mail.ParseAddress(emailText); err != nil {
			dialog.ShowError(fmt.Errorf("invalid email address"), ctx.App.ActiveWindow())
			return
		}

		log.Info("registration clicked",
			slog.String("email", emailText),
			slog.String("username", usernameText),
		)

		_, err := ctx.RPC.Register(context.TODO(), user.NewRegisterDTO(
			usernameText,
			emailText,
			passwordText,
		))

		if err != nil {
			dialog.ShowError(fmt.Errorf("registration failed: %v", err), ctx.App.ActiveWindow())
			return
		}

		usr, token, err := ctx.RPC.Login(context.TODO(), user.NewLoginDTO(
			emailText,
			passwordText,
		))

		if err != nil {
			dialog.ShowError(fmt.Errorf("something went wrong: %v", err), ctx.App.ActiveWindow())
			return
		}

		ctx.User = usr
		ctx.User.Token = token

		done <- struct{}{}
	}

	registerButton := widget.NewButtonWithIcon("Register", theme.ConfirmIcon(), submit)

	email.OnSubmitted = func(_ string) { submit() }
	username.OnSubmitted = func(_ string) { submit() }
	password.OnSubmitted = func(_ string) { submit() }

	login := widget.NewHyperlink("Back to Login", nil)
	login.OnTapped = back

	form := container.NewVBox(
		widget.NewLabelWithStyle("Create your account", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabel(""),
		widget.NewLabel("Email *"),
		email,
		widget.NewLabel("Username *"),
		username,
		widget.NewLabel("Password *"),
		password,
		registerButton,
		widget.NewLabel("Already have an account?"),
		login,
	)

	return form
}
