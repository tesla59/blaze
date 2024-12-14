package session

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/client"
	"log/slog"
	"sync"
)

type Session struct {
	ID      string
	Client1 *client.Client
	Client2 *client.Client
	Mu      *sync.Mutex
}

type Message struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewSession(id string, client1, client2 *client.Client) *Session {
	return &Session{
		ID:      id,
		Client1: client1,
		Client2: client2,
		Mu:      &sync.Mutex{},
	}
}

func (s *Session) HandleMessage(messageType int, messageByte []byte) {
	var message Message
	if err := json.Unmarshal(messageByte, &message); err != nil {
		slog.Error("Failed to unmarshal message", "error", err)
		return
	}
	if messageType != websocket.TextMessage {
		slog.Error("Invalid message type", "type", messageType)
		return
	}
	switch message.Type {
	case "message":
		s.handleTextMessage(message.Value)
	case "shuffle":
		s.handleShuffle()
	default:
	}
}

func (s *Session) handleTextMessage(message string) {
	respMessage := map[string]string{
		"type":  "message",
		"value": message,
	}
	if err := s.Client1.SendJSON(respMessage); err != nil {
		slog.Error("Failed to send message to client", "error", err)
	}
	if err := s.Client2.SendJSON(respMessage); err != nil {
		slog.Error("Failed to send message to client", "error", err)
	}
}

func (s *Session) handleShuffle() {
	clientCh := client.GetClientCh()
	if s.Client1 != nil {
		s.Client1.State = "waiting"
		clientCh <- s.Client1
	}
	if s.Client2 != nil {
		s.Client2.State = "waiting"
		clientCh <- s.Client2
		s.Client2 = nil // Reset client2
	}
}
