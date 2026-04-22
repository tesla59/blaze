package matchmaker

import (
	"sync"

	"github.com/tesla59/blaze/types"
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
			h.clients[client.Key()] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			// Remove the client from the hub
			h.Matchmaker.RemoveFromQueue(client)

			// Lock Hub after Matchmaker to prevent deadlock condition when queue is full
			h.mu.Lock()
			if _, ok := h.clients[client.Key()]; ok {
				delete(h.clients, client.Key())
				client.CloseChannels()
			}
			h.mu.Unlock()

			// Notify and re-enqueue peer if any
			if client.Peer != nil {
				peer := client.Peer
				peer.SafeSend(DisconnectedMessage())
				peer.State = types.Waiting
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
