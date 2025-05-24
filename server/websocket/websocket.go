package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/repository"
	"github.com/tesla59/blaze/service"
	"github.com/tesla59/blaze/types"
	"log/slog"
	"net/http"
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
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Identify the client
	_, message, err := conn.ReadMessage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	localClient, err := newClientFromMessage(message, conn, h.Hub)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slog.Debug("Client identified", "ID", localClient.ID)

	h.Hub.Register <- localClient

	go localClient.ReadPump()
	go localClient.WritePump()
}

// newClientFromMessage creates a new client from the initial message sent from the frontend.
func newClientFromMessage(message []byte, conn *websocket.Conn, h *matchmaker.Hub) (*matchmaker.Client, error) {
	slog.Debug("Received messageType", "message", string(message))

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

	client := matchmaker.NewClient(identityMessage.Client, "waiting", conn, h)

	return client, nil
}
