package server

import (
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/websocket/handler"
	"log/slog"
	"net/http"
)

type httpServer struct {
	cfg        config.Config
	mux        *http.ServeMux
	matchmaker *matchmaker.Matchmaker
}

func NewHTTPServer(cfg config.Config, mm *matchmaker.Matchmaker) Server {
	return &httpServer{
		cfg:        cfg,
		mux:        http.NewServeMux(),
		matchmaker: mm,
	}
}

func (s *httpServer) Start() error {
	s.registerHandlers()

	slog.Info("Starting Server on " + s.cfg.Server.Host + ":" + s.cfg.Server.Port)
	return http.ListenAndServe(s.cfg.Server.Host+":"+s.cfg.Server.Port, s.mux)
}

func (s *httpServer) registerHandlers() {
	s.mux.HandleFunc("/", homeHandler)
	s.mux.HandleFunc("/healthz", healthHandler)

	s.mux.HandleFunc("/ws", handler.NewWSHandler(s.matchmaker).Handle())

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./view/static"))))
}
