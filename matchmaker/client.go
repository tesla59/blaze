package matchmaker

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/types"
	"log/slog"
)

type Client struct {
	ID    string
	State string

	Session *Session
	Hub     *Hub

	Conn *websocket.Conn

	Send chan []byte

	Peer *Client
}

func NewClient(id, state string, conn *websocket.Conn, h *Hub) *Client {
	return &Client{
		ID:      id,
		State:   state,
		Session: nil,
		Hub:     h,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Peer:    nil,
	}
}

// HandleMessage handles incoming messages from the client
func (c *Client) HandleMessage(message []byte) {
	var messageType types.MessageType
	if err := json.Unmarshal(message, &messageType); err != nil {
		slog.Error("Failed to unmarshal message", "error", err)
		c.Send <- []byte("error: invalid message format")
		return
	}
	slog.Debug("Received message", "message", string(message))
	switch messageType.Type {
	case "join":
		slog.Debug("Client joined", "ID", c.ID)
		c.Hub.Matchmaker.Enqueue(c)
	}
}

func (c *Client) ReadPump() {
	// cleanup function to close the connection when the function exits
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			slog.Error("Failed to read message", "error", err)
			c.Send <- []byte("error: failed to read message")
			return
		}
		c.HandleMessage(message)
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
}
