package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/repository"
	"github.com/tesla59/blaze/service"
	"github.com/tesla59/blaze/types"
	"net/http"
	"slices"
)

type WSHandler struct {
	Upgrader      websocket.Upgrader
	Hub           *matchmaker.Hub
	DB            *pgxpool.Pool
	ClientService *service.ClientService
}

func NewWSHandler(hub *matchmaker.Hub, pool *pgxpool.Pool) *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				cfg := config.GetConfig()
				if cfg.Environment == "production" {
					origin := r.Header.Get("Origin")
					if slices.Contains(cfg.Server.AllowedOrigins, origin) {
						return true
					} else {
						log.Logger.Warn("Origin not allowed", "origin", origin, "allowedOrigins", cfg.Server.AllowedOrigins)
						return false
					}
				}
				return true
			},
		},
		Hub:           hub,
		DB:            pool,
		ClientService: service.NewClientService(repository.NewClientRepository(pool)),
	}
}

func (h *WSHandler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.websocketHandler(w, r)
	}
}

func (h *WSHandler) websocketHandler(w http.ResponseWriter, r *http.Request) {
	ctxLogger := log.Logger.With("handler", "websocketHandler", "path", r.URL.Path)
	ctx := log.Inject(r.Context(), ctxLogger)

	log.WithContext(ctx).Info("Handling websocket connection")

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.WithContext(ctx).Error("Failed to upgrade connection", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Identify the client
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.WithContext(ctx).Error("Failed to read initial message", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	localClient, err := newClientFromMessage(ctx, message, conn, h.Hub)
	if err != nil {
		log.WithContext(ctx).Error("Failed to create client from message", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.WithContext(ctx).Info("Client identified", "ID", localClient.ID)

	h.Hub.Register <- localClient

	ctx = log.Inject(ctx, log.Logger.With("clientID", localClient.ID))

	go localClient.ReadPump(ctx)
	go localClient.WritePump(ctx)
}

// newClientFromMessage creates a new client from the initial message sent from the frontend.
func newClientFromMessage(ctx context.Context, message []byte, conn *websocket.Conn, h *matchmaker.Hub) (*matchmaker.Client, error) {
	log.WithContext(ctx).Debug("Received messageType", "message", string(message))

	var identityMessage types.IdentityMessage

	if err := json.Unmarshal(message, &identityMessage); err != nil {
		return nil, fmt.Errorf("failed to unmarshal identity message: %w", err)
	}

	if identityMessage.Type != "identity" {
		return nil, fmt.Errorf("expected message type: identity, got: %s", identityMessage.Type)
	}

	if identityMessage.Client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	client := matchmaker.NewClient(identityMessage.Client, types.Connected, conn, h)

	return client, nil
}
