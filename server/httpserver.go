package server

import (
	"github.com/tesla59/blaze/config"
	"log/slog"
	"net/http"
)

type httpServer struct {
	cfg config.Config
	mux *http.ServeMux
}

func NewHTTPServer(cfg config.Config) Server {
	return &httpServer{
		cfg: cfg,
		mux: http.NewServeMux(),
	}
}

func (s *httpServer) Start() error {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("Hello, World!\n"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	slog.Info("Starting Server on " + s.cfg.Server.Host + ":" + s.cfg.Server.Port)
	return http.ListenAndServe(s.cfg.Server.Host+":"+s.cfg.Server.Port, s.mux)
}
