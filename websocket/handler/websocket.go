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

	localClient, err := getClientFromMessage(message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	localClient.Conn = conn
	localClient.RandomizeTarget()
	slog.Debug("Client identified", "ID", localClient.ID, "ConnectedTo", localClient.ConnectedTo)

	clientMap := client.GetClientMap()
	clientMap.AddClient(&localClient)
	defer clientMap.RemoveClient(localClient.ID)

	slog.Debug("Client connected", "ID", localClient.ID)

	// Remove Later
	resp := map[string]string{
		"type":   "identity",
		"id":     localClient.ID,
		"target": localClient.ConnectedTo,
	}
	respBytes, _ := json.Marshal(resp)
	if err := conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
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

func getClientFromMessage(message []byte) (client.Client, error) {
	identityMessage := make(map[string]string)
	slog.Debug("Received message", "message", string(message))

	if err := json.Unmarshal(message, &identityMessage); err != nil {
		return client.Client{}, err
	}

	if _, ok := identityMessage["id"]; !ok {
		return client.Client{}, fmt.Errorf("ID not found")
	}

	id := identityMessage["id"]

	return client.Client{
		ID:          id,
		Conn:        nil,
		ConnectedTo: "",
	}, nil
}

func handleClientMessage(c *client.Client, messageType int, message []byte) {
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
			"to":      c.ConnectedTo,
		}
		messageByte, err := json.Marshal(messageMap)
		if err != nil {
			slog.Warn("Error marshalling message", "ID", c.ID, "error", err)
			return
		}
		c.SendMessage(c.ConnectedTo, messageType, messageByte)
	case "shuffle":
		c.RandomizeTarget()
		resp := map[string]string{
			"type":   "identity",
			"id":     c.ID,
			"target": c.ConnectedTo,
		}
		respBytes, _ := json.Marshal(resp)
		if err := c.Conn.WriteMessage(websocket.TextMessage, respBytes); err != nil {
			slog.Warn("Error sending message", "ID", c.ID, "error", err)
			return
		}
		slog.Debug("Client shuffled", "ID", c.ID, "ConnectedTo", c.ConnectedTo)
	default:
		slog.Warn("Unknown message type", "ID", c.ID, "type", messageJSON["type"])
	}
}
