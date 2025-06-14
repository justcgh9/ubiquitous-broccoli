package loginform

import (
	"fmt"
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func NewLoginForm(log *slog.Logger) *fyne.Container {
	email := widget.NewEntry()
	email.SetPlaceHolder("Email")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Password")

	loginButton := widget.NewButtonWithIcon("Log In", theme.ConfirmIcon(), func() {
		log.Info(fmt.Sprintf("login clicked: %s %s", email.Text, password.Text))
	})

	forgot := widget.NewHyperlink("Forgot your password?", nil)
	register := widget.NewHyperlink("Register", nil)

	return container.NewVBox(
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
}