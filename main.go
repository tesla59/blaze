package main

import (
	"context"
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/database"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/server"
	"log/slog"
)

func main() {
	cfg := config.GetConfig()
	log.Init()

	ctx := context.Background()
	db, err := database.GetPool(ctx)
	if err != nil {
		log.Logger.Error("error connecting to database", "error", err.Error())
		return
	}

	matchMaker := matchmaker.NewMatchmaker(100)
	hub := matchmaker.NewHub(matchMaker)
	go hub.Run()

	httpServer := server.NewHTTPServer(cfg, hub, db)
	slog.SetLogLoggerLevel(slog.LevelDebug)

	err = httpServer.Start()
	if err != nil {
		slog.Error("error starting server", "error", err.Error())
	}
}
