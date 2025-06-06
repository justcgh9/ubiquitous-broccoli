package main

import (
	"log/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
    log := slog.New(&slog.TextHandler{})

    a := app.New()
    loginWindow := a.NewWindow("ForkCord")
    _ = loginWindow

    bg := canvas.NewImageFromFile("media/background-login.png")
    bg.FillMode = canvas.ImageFillStretch

    email := widget.NewEntry()
    email.SetPlaceHolder("Email")

    password := widget.NewPasswordEntry()
    password.SetPlaceHolder("Password")

    loginButton := widget.NewButtonWithIcon("Log In", theme.ConfirmIcon(), func() {
        log.Info("Login clicked: ", email.Text, " ", password.Text)
    })

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

    card := container.NewPadded(container.NewVBox(form))
    cardBG := canvas.NewRectangle(theme.Color(theme.ColorNameBackground))
    cardBG.SetMinSize(fyne.NewSize(320, 400))
    loginCard := container.NewStack(cardBG, card)

    content := container.NewStack(
        bg,
        container.NewCenter(loginCard),
    )

    loginWindow.SetContent(content)
    loginWindow.Resize(fyne.NewSize(900, 600))
    loginWindow.Show()
    a.Run()
}

