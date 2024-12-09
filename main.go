package main

import (
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/server"
	"log/slog"
)

func main() {
	cfg := config.GetConfig()
	httpServer := server.NewHTTPServer(*cfg)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := httpServer.Start()
	if err != nil {
		slog.Error("error starting server", "error", err.Error())
	}
}
