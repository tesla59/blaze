package main

import (
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/server"
	"log/slog"
)

func main() {
	cfg := config.GetConfig()

	matchMaker := matchmaker.NewMatchmaker(100)
	hub := matchmaker.NewHub(matchMaker)
	go hub.Run()

	httpServer := server.NewHTTPServer(cfg, hub)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err := httpServer.Start()
	if err != nil {
		slog.Error("error starting server", "error", err.Error())
	}
}
