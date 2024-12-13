package session

import "github.com/tesla59/blaze/client"

type Session struct {
	ID      string
	Client1 client.Client
	Client2 client.Client
}

func NewSession(id string, client1, client2 client.Client) *Session {
	return &Session{
		ID:      id,
		Client1: client1,
		Client2: client2,
	}
}
