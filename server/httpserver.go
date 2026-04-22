package server

import (
	"context"
	"crypto/tls"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/server/client"
	serveMatchmaker "github.com/tesla59/blaze/server/matchmaker"
	"github.com/tesla59/blaze/server/websocket"
	"github.com/tesla59/blaze/types"
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
		Handler: corsMiddleware(mux),
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
	log.Logger.Info("Starting Server on " + s.cfg.Server.Host + ":" + s.cfg.Server.Port)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Logger.Info("Shutting down server...")

		// wait for 10 seconds for all websocket connections to exist gracefully
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			log.Logger.Error("Server shutdown error", "error", err)
		}
	}()

	var err error
	if s.cfg.Server.SSL.Enabled {
		err = s.server.ListenAndServeTLS(s.cfg.Server.SSL.CertFile, s.cfg.Server.SSL.KeyFile)
	} else {
		err = s.server.ListenAndServe()
	}

	if errors.Is(err, http.ErrServerClosed) {
		log.Logger.Info("Server stopped gracefully")
		return nil
	}
	return err
}

func (s *httpServer) registerHandlers() {
	handlerMap := map[string]http.HandlerFunc{
		"/":                          homeHandler,
		"/healthz":                   healthHandler,
		"/ws":                        websocket.NewWSHandler(s.hub, s.db).Handle(),
		"POST /api/v1/client":        client.NewClientHandler(s.db).Handle("POST"),
		"GET /api/v1/client":         client.NewClientHandler(s.db).Handle("GET"),
		"POST /api/v1/client/verify": client.NewClientHandler(s.db).Handle("POST"),
	}
	if s.cfg.Environment != "production" {
		handlerMap["/queue"] = serveMatchmaker.NewQueueStateHandler(s.hub.Matchmaker).Handle()
	}
	for path, handler := range handlerMap {
		s.mux.HandleFunc(path, handler)
	}

	s.mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./view/static"))))
}

// corsMiddleware is a middleware function that adds CORS headers to the response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowedOrigins := config.GetConfig().Server.AllowedOrigins

		// Enforce CORS in production environment only. Skip on development
		if config.GetConfig().Environment == "production" {
			if origin == "" || !slices.Contains(allowedOrigins, origin) {
				log.Logger.Warn("Origin not allowed", "origin", origin, "allowedOrigins", allowedOrigins)
				http.Error(w, "Origin not allowed", http.StatusForbidden)
				return
			}
		}

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin") // important for caching proxies
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
