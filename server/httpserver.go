package server

import (
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/server/websocket"
	"github.com/tesla59/blaze/types"
	"log/slog"
	"net/http"
)

type httpServer struct {
	cfg *types.Config
	mux *http.ServeMux
	hub *matchmaker.Hub
}

func NewHTTPServer(cfg *types.Config, hub *matchmaker.Hub) Server {
	return &httpServer{
		cfg: cfg,
		mux: http.NewServeMux(),
		hub: hub,
	}
}

func (s *httpServer) Start() error {
	s.registerHandlers()

	slog.Info("Starting Server on " + s.cfg.Server.Host + ":" + s.cfg.Server.Port)
	return http.ListenAndServe(s.cfg.Server.Host+":"+s.cfg.Server.Port, s.mux)
}

func (s *httpServer) registerHandlers() {
	handlerMap := map[string]http.HandlerFunc{
		"/":        homeHandler,
		"/healthz": healthHandler,
		"/ws":      websocket.NewWSHandler(s.hub).Handle(),
	}
	for path, handler := range handlerMap {
		s.mux.HandleFunc(path, handler)
	}

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./view/static"))))
}
