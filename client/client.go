package client

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Client struct {
	ID    string
	State string
	Conn  *websocket.Conn
	Mutex sync.RWMutex
}

func NewClient(id, state string, conn *websocket.Conn) Client {
	mut := sync.RWMutex{}
	return Client{
		ID:    id,
		State: state,
		Conn:  conn,
		Mutex: mut,
	}
}

func (c *Client) SendMessage(messageType int, message []byte) error {
	return c.Conn.WriteMessage(messageType, message)
}

func (c *Client) SendJSON(message map[string]string) error {
	return c.Conn.WriteJSON(message)
}
