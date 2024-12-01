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

	clientMap := client.GetClientMap()
	clientMap.AddClient(localClient)
	defer clientMap.RemoveClient(localClient.ID)

	slog.Debug("Client connected", "ID", localClient.ID)

	// Remove Later
	conn.WriteMessage(websocket.TextMessage, []byte("Connected to server"))

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send the message to the target client
		localClient.SendMessage(localClient.ConnectedTo, messageType, string(message)+" from "+localClient.ID)
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
