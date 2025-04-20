package matchmaker

import (
	"log/slog"
	"sync"
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
	slog.Debug("Enqueueing client", "ID", c.ID)
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the client is already in the queue
	for len(m.queue) > 0 {
		peer := m.queue[0]
		if peer.ID == c.ID {
			continue
		}
		if peer.State == "queued" {
			m.queue = m.queue[1:]
			m.matchPair(c, peer)
			return
		}
		// drop the peer
		m.queue = m.queue[1:]
	}

	// Add the client to the queue
	c.State = "queued"
	m.queue = append(m.queue, c)
}

// matchPair creates a session, ties the two clients together, and notifies them
func (m *Matchmaker) matchPair(a, b *Client) {
	slog.Debug("Matched pair", "a", a.ID, "b", b.ID)
	// Create a session and tie the two clients together
	session := NewSession(a, b)

	a.Session = session
	a.Peer = b
	a.State = "matched"

	b.Session = session
	b.Peer = a
	b.State = "matched"

	a.Send <- []byte("matched")
	b.Send <- []byte("matched")
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
