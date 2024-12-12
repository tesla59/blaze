package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/client"
	"log/slog"
	"net/http"
)

type WSHandler struct {
	Upgrader websocket.Upgrader
}

func NewWSHandler() *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
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
	defer conn.Close()

	// Identify the client
	_, message, err := conn.ReadMessage()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	localClient, err := getClientFromMessage(message, conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Debug("Client identified", "ID", localClient.ID)
	client.ClientCh <- localClient

	// Remove Later
	resp := map[string]string{
		"type": "identity",
		"id":   localClient.ID,
	}
	if err := localClient.SendJSON(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		handleClientMessage(&localClient, messageType, message)
	}
}

// getClientFromMessage extracts the client ID from the initial message sent from frontend and returns a new client
func getClientFromMessage(message []byte, conn *websocket.Conn) (client.Client, error) {
	identityMessage := make(map[string]string)
	slog.Debug("Received message", "message", string(message))

	if err := json.Unmarshal(message, &identityMessage); err != nil {
		return client.Client{}, err
	}

	if _, ok := identityMessage["id"]; !ok {
		return client.Client{}, fmt.Errorf("ID not found")
	}

	id := identityMessage["id"]

	return client.NewClient(id, "waiting", conn), nil
}

// handleClientMessage handles JSON Message from the client
func handleClientMessage(c *client.Client, _ int, message []byte) {
	var messageJSON map[string]string
	if err := json.Unmarshal(message, &messageJSON); err != nil {
		slog.Warn("Error unmarshalling message", "ID", c.ID, "error", err)
		return
	}

	switch messageJSON["type"] {
	case "message":
		messageMap := map[string]string{
			"type":    "message",
			"message": messageJSON["message"],
			"from":    c.ID,
		}
		if err := c.SendJSON(messageMap); err != nil {
			slog.Warn("Error sending message", "ID", c.ID, "error", err)
			return
		}
	case "shuffle":
		resp := map[string]string{
			"type": "identity",
			"id":   c.ID,
		}
		respBytes, _ := json.Marshal(resp)
		if err := c.Conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			slog.Warn("Error sending message", "ID", c.ID, "error", err)
			return
		}
	default:
		slog.Warn("Unknown message type", "ID", c.ID, "type", messageJSON["type"])
	}
}
