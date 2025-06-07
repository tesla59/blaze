package matchmaker

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/types"
	"time"
)

const (
	// writeWait is the time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// pongWait is the time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 1024 * 1024
)

type Client struct {
	*models.Client
	State types.State

	Session *Session
	Hub     *Hub

	Conn *websocket.Conn

	Send chan []byte

	Peer *Client
}

func NewClient(c *models.Client, state types.State, conn *websocket.Conn, h *Hub) *Client {
	return &Client{
		Client:  c,
		State:   state,
		Session: nil,
		Hub:     h,
		Conn:    conn,
		Send:    make(chan []byte, 256),
		Peer:    nil,
	}
}

// HandleMessage handles incoming messages from the client
func (c *Client) HandleMessage(ctx context.Context, message []byte) {
	var messageType types.MessageType
	if err := json.Unmarshal(message, &messageType); err != nil {
		log.WithContext(ctx).Error("Failed to unmarshal received message", "message", string(message), "error", err)
		c.Send <- ErrorByte(err)
		return
	}
	log.WithContext(ctx).Debug("Received message", "message", string(message))
	switch messageType.Type {
	case "join":
		log.WithContext(ctx).Debug("Client joined")
		c.Hub.Matchmaker.Enqueue(c)
	case "message":
		var peerMessage types.Message
		if err := json.Unmarshal(message, &peerMessage); err != nil {
			log.WithContext(ctx).Error("Failed to unmarshal peer message", "error", err)
			c.Send <- ErrorByte(err)
			return
		}
		log.WithContext(ctx).Info("Forwarding message", "message", peerMessage.Message)
		if c.Peer == nil {
			log.WithContext(ctx).Error("No peer to send message to", "ID", c.ID)
			c.Send <- ErrorByte(errors.New("no peer connected"))
			return
		}
		c.Peer.Send <- message
	case "end":
		log.WithContext(ctx).Info("Client ended chat session", "ID", c.ID)
		a, b := c.Peer, c
		if a != nil {
			log.WithContext(ctx).Info("Disconnecting peer", "ID", a.ID)
			a.Session = nil
			a.Peer = nil
			a.State = types.Waiting
			a.Send <- DisconnectedMessage()
			c.Hub.Matchmaker.Enqueue(a)
		}
		if b != nil {
			b.Peer = nil
			b.State = types.Connected
			b.Session = nil
			b.Send <- DisconnectedMessage()
		}
	case "disconnect":
		log.WithContext(ctx).Info("Client disconnected")
		c.Hub.Unregister <- c
	case "rematch":
		log.WithContext(ctx).Info("Client rematch", "ID", c.ID)
		a, b := c.Peer, c
		if a != nil {
			a.Session = nil
			a.Peer = nil
			a.State = types.Waiting
			a.Send <- DisconnectedMessage()
			c.Hub.Matchmaker.Enqueue(a)
		}
		if b != nil {
			b.Peer = nil
			b.State = types.Waiting
			b.Send <- DisconnectedMessage()
			c.Hub.Matchmaker.Enqueue(b)
		}
	case "sdp-offer", "sdp-answer", "ice-candidate":
		if c.Peer != nil {
			log.WithContext(ctx).Debug("Forwarding message to peer", "type", messageType.Type)
			c.Peer.Send <- message
		} else {
			log.WithContext(ctx).Error("No peer to forward message to")
			c.Send <- ErrorByte(errors.New("no peer connected"))
		}
	default:
		log.WithContext(ctx).Error("Unknown message type", "type", messageType.Type)
	}
}

func (c *Client) ReadPump(ctx context.Context) {
	// cleanup function to close the connection when the function exits
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		log.WithContext(ctx).Debug("Received pong", "ID", c.ID)
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			var closeErr *websocket.CloseError
			if errors.As(err, &closeErr) && closeErr.Code == websocket.CloseGoingAway {
				log.WithContext(ctx).Info("Client connection closed by peer", "ID", c.ID)
			} else {
				log.WithContext(ctx).Error("Failed to read message", "ID", c.ID, "error", err)
			}
			return
		}
		c.HandleMessage(ctx, message)
	}
}

func (c *Client) WritePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.WithContext(ctx).Error("Failed to get next writer", "error", err)
				return
			}
			if _, err := w.Write(message); err != nil {
				log.WithContext(ctx).Error("Failed to write message", "error", err)
				return
			}
			if err := w.Close(); err != nil {
				log.WithContext(ctx).Error("Failed to close writer", "error", err)
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.WithContext(ctx).Error("Failed to write ping message", "error", err)
				return
			}
		}
	}
}
