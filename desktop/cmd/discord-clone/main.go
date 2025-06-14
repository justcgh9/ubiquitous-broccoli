package main

import (
	"log/slog"

	"fyne.io/fyne/v2/app"
    myApp "github.com/justcgh9/discord-clone/desktop/internal/app"
)

func main() {
    log := slog.New(&slog.TextHandler{})

    a := myApp.New(
        log,
        app.New(),
    )

    a.Run()
}