package server

import (
	"crypto/tls"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/server/client"
	serveMatchmaker "github.com/tesla59/blaze/server/matchmaker"
	"github.com/tesla59/blaze/server/websocket"
	"github.com/tesla59/blaze/types"
	"log/slog"
	"net/http"
)

type httpServer struct {
	cfg    *types.Config
	mux    *http.ServeMux
	server *http.Server
	hub    *matchmaker.Hub
	db     *pgxpool.Pool
}

func NewHTTPServer(cfg *types.Config, hub *matchmaker.Hub, pool *pgxpool.Pool) Server {
	mux := http.NewServeMux()
	serv := &http.Server{
		Addr:    cfg.Server.Host + ":" + cfg.Server.Port,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS13,
		},
	}
	return &httpServer{
		cfg:    cfg,
		server: serv,
		mux:    mux,
		hub:    hub,
		db:     pool,
	}
}

func (s *httpServer) Start() error {
	s.registerHandlers()
	slog.Info("Starting Server on " + s.cfg.Server.Host + ":" + s.cfg.Server.Port)

	if s.cfg.Server.SSL.Enabled {
		return s.server.ListenAndServeTLS(s.cfg.Server.SSL.CertFile, s.cfg.Server.SSL.KeyFile)
	} else {
		return s.server.ListenAndServe()
	}
}

func (s *httpServer) registerHandlers() {
	handlerMap := map[string]http.HandlerFunc{
		"/":                          homeHandler,
		"/healthz":                   healthHandler,
		"/ws":                        websocket.NewWSHandler(s.hub, s.db).Handle(),
		"/queue":                     serveMatchmaker.NewQueueStateHandler(s.hub.Matchmaker).Handle(),
		"POST /api/v1/client":        client.NewClientHandler(s.db).Handle("POST"),
		"GET /api/v1/client":         client.NewClientHandler(s.db).Handle("GET"),
		"POST /api/v1/client/verify": client.NewClientHandler(s.db).Handle("POST"),
	}
	for path, handler := range handlerMap {
		s.mux.HandleFunc(path, handler)
	}

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./view/static"))))
}
