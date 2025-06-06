package matchmaker

import (
	"strconv"
	"sync"

	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/types"
)

type Matchmaker struct {
	mu    sync.Mutex
	queue []*Client
}

func NewMatchmaker(queueSize int) *Matchmaker {
	return &Matchmaker{
		queue: make([]*Client, 0, queueSize),
		mu:    sync.Mutex{},
	}
}

func (m *Matchmaker) Enqueue(c *Client) {
	log.Logger.Info("Enqueueing client", "ID", c.ID)
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the client is already in the queue
	for _, queued := range m.queue {
		if queued.ID == c.ID {
			log.Logger.Info("Client already in queue", "ID", c.ID)
			return
		}
	}

	// Match with the first available peer
	for i, peer := range m.queue {
		if peer.State == types.Waiting && peer.ID != c.ID {
			// Remove the peer from the queue
			m.queue = append(m.queue[:i], m.queue[i+1:]...)
			m.matchPair(c, peer)
			return
		}
	}

	// Add the client to the queue
	c.State = types.Waiting
	m.queue = append(m.queue, c)
}

// matchPair creates a session, ties the two clients together, and notifies them
func (m *Matchmaker) matchPair(a, b *Client) {
	log.Logger.Info("Matched pair", "a", a.ID, "b", b.ID)
	// Create a session and tie the two clients together
	session := NewSession(a, b)

	a.Session = session
	a.Peer = b
	a.State = types.Matched

	b.Session = session
	b.Peer = a
	b.State = types.Matched

	a.Send <- MatchedMessage(b)
	b.Send <- MatchedMessage(a)
}

// RemoveFromQueue removes c if it’s still waiting
func (m *Matchmaker) RemoveFromQueue(c *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, queued := range m.queue {
		if queued.ID == c.ID {
			// drop it out
			m.queue = append(m.queue[:i], m.queue[i+1:]...)
			break
		}
	}
}

func (m *Matchmaker) GetQueueState() []map[string]string {
	m.mu.Lock()
	defer m.mu.Unlock()

	queueState := make([]map[string]string, len(m.queue))
	for i, client := range m.queue {
		queueState[i] = map[string]string{
			"ID":    strconv.Itoa(client.ID),
			"State": client.State.String(),
		}
	}
	return queueState
}
