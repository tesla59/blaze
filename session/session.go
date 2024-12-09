package session

import "github.com/tesla59/blaze/client"

type Session struct {
	ID      string
	Client1 client.Client
	Client2 client.Client
}
