package matchmaker

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tesla59/blaze/client"
	"github.com/tesla59/blaze/session"
	"sync"
	"time"
)

type Matchmaker struct {
	sessions map[string]*session.Session
	mu       sync.Mutex
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		sessions: make(map[string]*session.Session),
	}
}

func (m *Matchmaker) Start() {
	for {
		newClient := <-client.ClientCh
		m.mu.Lock()

		var matchedSession *session.Session
		for _, session := range m.sessions {
			if session.Client2.State == "waiting" {
				matchedSession = session
				break
			}
		}

		if matchedSession != nil { // Match found
			matchedSession.Client2 = newClient
			newClient.State = "matched"

			matchedSession.Client1.SendMessage(websocket.TextMessage, []byte("Matched with client: "+newClient.ID))
			newClient.SendMessage(websocket.TextMessage, []byte("Matched with client: "+matchedSession.Client1.ID))

			fmt.Printf("Client %s matched with Client %s in Session %s\n",
				matchedSession.Client1.ID, newClient.ID, matchedSession.ID)
		} else { // No match found, create new session
			sessionID := generateSessionID()
			newSession := &session.Session{
				ID:      sessionID,
				Client1: newClient,
			}
			m.sessions[sessionID] = newSession
			newClient.State = "waiting"

			fmt.Printf("New session created with ID: %s for Client %s\n", sessionID, newClient.ID)
		}

		m.mu.Unlock()
	}
}

func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}
