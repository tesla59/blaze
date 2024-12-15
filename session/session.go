package session

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/client"
	"github.com/tesla59/blaze/types"
	"log/slog"
	"sync"
)

type Session struct {
	ID      string
	Client1 *client.Client
	Client2 *client.Client
	Mu      *sync.Mutex
}

func NewSession(id string, client1, client2 *client.Client) *Session {
	return &Session{
		ID:      id,
		Client1: client1,
		Client2: client2,
		Mu:      &sync.Mutex{},
	}
}

func (s *Session) HandleMessage(messageType int, messageByte []byte, clientID string) {
	var message types.Message
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
		s.handleTextMessage(clientID, message.Value)
	case "shuffle":
		s.handleShuffle()
	case "disconnect":
		s.handleDisconnect(clientID)
	default:
	}
}

func (s *Session) handleTextMessage(clientID, message string) {
	respMessage := map[string]string{
		"type":  "message",
		"value": message,
	}
	if clientID == s.Client1.ID {
		if err := s.Client1.SendJSON(respMessage); err != nil {
			slog.Error("Failed to send message to client 1", "client", s.Client1.ID, "error", err)
		}
	} else {
		if err := s.Client2.SendJSON(respMessage); err != nil {
			slog.Error("Failed to send message to client 2", "client", s.Client2.ID, "error", err)
		}
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

func (s *Session) handleDisconnect(clientID string) {
	clientCh := client.GetClientCh()
	resp := map[string]string{
		"type":  "rematch",
		"value": "You have disconnected. Waiting for a new match...",
	}

	if s.Client1.ID == clientID {
		nullClient := client.NullClient(s.ID, s.Client1.Conn)
		s.Client1 = &nullClient

		if err := s.Client2.SendJSON(resp); err != nil {
			slog.Error("Failed to send rematch message to client 2", "error", err)
		}
		clientCh <- s.Client2
	} else {
		nullClient := client.NullClient(s.ID, s.Client2.Conn)
		s.Client2 = &nullClient

		if err := s.Client1.SendJSON(resp); err != nil {
			slog.Error("Failed to send rematch message to client 1", "error", err)
		}
		clientCh <- s.Client1
	}
}
