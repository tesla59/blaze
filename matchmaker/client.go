package matchmaker

import (
	"encoding/json"
	"errors"
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
		slog.Error("Failed to unmarshal message", "ID", c.ID, "error", err)
		c.Send <- ErrorByte(err)
		return
	}
	slog.Debug("Received message", "message", string(message))
	switch messageType.Type {
	case "join":
		slog.Debug("Client joined", "ID", c.ID)
		c.Hub.Matchmaker.Enqueue(c)
	case "message":
		var peerMessage types.Message
		if err := json.Unmarshal(message, &peerMessage); err != nil {
			slog.Error("Failed to unmarshal peer message", "ID", c.ID, "error", err)
			c.Send <- ErrorByte(err)
			return
		}
		slog.Debug("Client message", "ID", c.ID, "message", peerMessage.Message)
		c.Peer.Send <- message
	case "disconnect":
		slog.Debug("Client disconnected", "ID", c.ID)
		c.Hub.Unregister <- c
	case "rematch":
		slog.Debug("Client rematch", "ID", c.ID, "peerID", c.Peer.ID)
		a, b := c.Peer, c
		if a != nil {
			a.Session = nil
			a.Peer = nil
			a.State = "queued"
			a.Send <- DisconnectedMessage()
			c.Hub.Matchmaker.Enqueue(a)
		}
		if b != nil {
			b.Peer = nil
			b.State = "queued"
			b.Send <- DisconnectedMessage()
			c.Hub.Matchmaker.Enqueue(b)
		}
	case "sdp-offer", "sdp-answer", "ice-candidate":
		if c.Peer != nil {
			slog.Debug("Forwarding message to peer", "peerID", c.Peer.ID, "type", messageType.Type)
			c.Peer.Send <- message
		} else {
			slog.Error("No peer to forward message to", "ID", c.ID)
			c.Send <- ErrorByte(errors.New("no peer connected"))
		}
	default:
		slog.Error("Unknown message type", "ID", c.ID, "type", messageType.Type)
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
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) && closeErr.Code == websocket.CloseGoingAway {
				slog.Info("Client connection closed by peer", "ID", c.ID)
			} else {
				slog.Error("Failed to read message", "ID", c.ID, "error", err)
			}
			return
		}
		c.HandleMessage(message)
	}
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			slog.Error("Failed to write message", "ID", c.ID, "error", err)
			break
		}
	}
}
