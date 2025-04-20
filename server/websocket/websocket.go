package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/matchmaker"
	"github.com/tesla59/blaze/types"
	"log/slog"
	"net/http"
)

type WSHandler struct {
	Upgrader websocket.Upgrader
	Hub      *matchmaker.Hub
}

func NewWSHandler(hub *matchmaker.Hub) *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Hub: hub,
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

// getClientFromMessage extracts the client ID from the initial message sent from frontend and returns a new client
func newClientFromMessage(message []byte, conn *websocket.Conn, h *matchmaker.Hub) (*matchmaker.Client, error) {
	var messageType types.MessageType
	slog.Debug("Received messageType", "message", string(message))

	if err := json.Unmarshal(message, &messageType); err != nil {
		return nil, err
	}

	if messageType.Type != "identity" {
		return nil, fmt.Errorf("expected message type: identity, got: %s", messageType.Type)
	}

	var identityMessage types.IdentityMessage

	if err := json.Unmarshal(message, &identityMessage); err != nil {
		return nil, err
	}

	return matchmaker.NewClient(identityMessage.ClientID, "waiting", conn, h), nil
}
