package main

import (
	"log/slog"
	"os"

	"github.com/justcgh9/discord-clone-users/internal/lib/logger/handlers/pretty"
)

func main() {
	h := pretty.NewPrettyHandler(os.Stdout, pretty.PrettyHandlerOptions{
		ShowCaller: true,
		IndentJSON: true,
	})
	
	logger := slog.New(h)
	logger.Info("hello world", slog.String("foo", "bar"))
}