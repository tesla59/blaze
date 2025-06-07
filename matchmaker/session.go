package matchmaker

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type session struct {
	ID      string
	Client1 *Client
	Client2 *Client
	Mu      *sync.Mutex
}

func newSession(client1, client2 *Client) *session {
	return &session{
		ID:      generateSessionID(),
		Client1: client1,
		Client2: client2,
		Mu:      &sync.Mutex{},
	}
}

func generateSessionID() string {
	return fmt.Sprintf("session-%d-%d", time.Now().UnixNano(), rand.Int())
}
