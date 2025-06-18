package main

import (
	"log/slog"
	"os"

	"fyne.io/fyne/v2/app"
	myApp "github.com/justcgh9/discord-clone/desktop/internal/app"
)

func main() {
    log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))

    a := myApp.Run(
        log,
        app.New(),
        "localhost:44044",
    )

    _ = a
}