package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	once sync.Once
)

type Client struct {
	ID        string
	State     string
	Conn      *websocket.Conn
	Mutex     *sync.RWMutex
	SessionID string
}

func NewClient(id, state, sessionID string, conn *websocket.Conn) *Client {
	mut := sync.RWMutex{}
	return &Client{
		ID:        id,
		State:     state,
		Conn:      conn,
		Mutex:     &mut,
		SessionID: sessionID,
	}
}

func NullClient(sessionID string, conn *websocket.Conn) Client {
	return Client{
		ID:        "null",
		State:     "waiting",
		Conn:      conn,
		Mutex:     &sync.RWMutex{},
		SessionID: sessionID,
	}
}

// SendMessage sends a message to the client. Deprecated: Use SendJSON instead.
func (c *Client) SendMessage(messageType int, message []byte) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Conn.WriteMessage(messageType, message)
}

// SendJSON sends a JSON message to the client. Prefer this over SendMessage.
func (c *Client) SendJSON(message map[string]string) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	return c.Conn.WriteJSON(message)
}
