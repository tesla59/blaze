package matchmaker

import (
	"strconv"
	"sync"
)

type Hub struct {
	clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
	Matchmaker *Matchmaker
	mu         sync.RWMutex
}

func NewHub(matchmaker *Matchmaker) *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Matchmaker: matchmaker,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.clients[strconv.Itoa(client.ID)] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			// Remove the client from the hub
			h.Matchmaker.RemoveFromQueue(client)
			if _, ok := h.clients[strconv.Itoa(client.ID)]; ok {
				delete(h.clients, strconv.Itoa(client.ID))
				close(client.Send)
			}
			h.mu.Unlock()

			// Notify and re-enqueue peer if any
			if client.Peer != nil {
				peer := client.Peer
				peer.Send <- DisconnectedMessage()
				peer.State = "queued"
				peer.Peer = nil
				h.Matchmaker.Enqueue(peer)
			}
		}
	}
}

func (h *Hub) GetClientByID(id string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clients[id]
	return client, ok
}
