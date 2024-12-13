package matchmaker

import (
	"fmt"
	"github.com/tesla59/blaze/client"
	"github.com/tesla59/blaze/session"
	"log/slog"
	"sync"
	"time"
)

type Matchmaker struct {
	Sessions map[string]*session.Session
	ClientCh chan *client.Client
	Mu       sync.Mutex
}

func NewMatchmaker() *Matchmaker {
	return &Matchmaker{
		Sessions: make(map[string]*session.Session),
		ClientCh: make(chan *client.Client),
		Mu:       sync.Mutex{},
	}
}

func (m *Matchmaker) Start() {
	for {
		newClient := <-m.ClientCh
		m.Mu.Lock()

		var matchedSession *session.Session
		var firstClientEmpty bool
		for _, session := range m.Sessions {
			if session.Client2.State == "waiting" || session.Client1.State == "waiting" {
				if session.Client1.State == "waiting" {
					firstClientEmpty = false
				} else {
					firstClientEmpty = true
				}
				matchedSession = session
				break
			}
		}

		if matchedSession != nil { // Match found
			if firstClientEmpty {
				matchedSession.Client1 = newClient
			} else {
				matchedSession.Client2 = newClient
			}
			matchedSession.Client1.State = "matched"
			matchedSession.Client2.State = "matched"
			matchedSession.Client1.SessionID = matchedSession.ID
			matchedSession.Client2.SessionID = matchedSession.ID

			respMessage := map[string]string{
				"type":  "message",
				"value": "Matched with client: " + matchedSession.Client1.ID,
			}
			matchedSession.Client2.SendJSON(respMessage)
			respMessage["message"] = "Matched with client: " + matchedSession.Client2.ID
			matchedSession.Client1.SendJSON(respMessage)

			fmt.Printf("Client %s matched with Client %s in Session %s\n",
				matchedSession.Client1.ID, newClient.ID, matchedSession.ID)
		} else { // No match found, create new session
			sessionID := generateSessionID()
			newClient.State = "waiting"
			newClient.SessionID = sessionID

			newSession := session.NewSession(sessionID, newClient, client.NewClient("null", "waiting", sessionID, nil))
			m.Sessions[sessionID] = newSession
			slog.Debug("New session created with ID: %s for Client %s\n", sessionID, newClient.ID)
		}

		m.Mu.Unlock()
	}
}

func generateSessionID() string {
	return fmt.Sprintf("session-%d", time.Now().UnixNano())
}
