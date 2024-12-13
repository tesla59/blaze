package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/client"
	"github.com/tesla59/blaze/matchmaker"
	"log/slog"
	"net/http"
	"time"
)

type WSHandler struct {
	Upgrader   websocket.Upgrader
	MatchMaker *matchmaker.Matchmaker
}

type Message struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewWSHandler(mm *matchmaker.Matchmaker) *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		MatchMaker: mm,
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
	h.MatchMaker.ClientCh <- localClient
	time.Sleep(1 * time.Second)
	session, ok := h.MatchMaker.Sessions[localClient.SessionID]
	if !ok {
		slog.Error("Session not found", "sessionID", localClient.SessionID)
		return
	}

	for {
		messageType, messageByte, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.HandleMessage(messageType, messageByte)
	}
}

// getClientFromMessage extracts the client ID from the initial message sent from frontend and returns a new client
func getClientFromMessage(message []byte, conn *websocket.Conn) (*client.Client, error) {
	var identityMessage Message
	slog.Debug("Received message", "message", string(message))

	if err := json.Unmarshal(message, &identityMessage); err != nil {
		return nil, err
	}

	if identityMessage.Type != "identity" {
		return nil, fmt.Errorf("expected message type: identity, got: %s", identityMessage.Type)
	}

	return client.NewClient(identityMessage.Value, "waiting", "", conn), nil
}
