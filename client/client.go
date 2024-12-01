package client

import (
	"github.com/gorilla/websocket"
	"log/slog"
	"sync"
)

type Map struct {
	Map   map[string]*Client
	Mutex sync.RWMutex
}

type Client struct {
	ID          string
	ConnectedTo string
	Conn        *websocket.Conn
	Mutex       sync.RWMutex
}

var ClientMap *Map
var once = sync.Once{}

func GetClientMap() *Map {
	once.Do(func() {
		slog.Debug("Creating new client map")
		ClientMap = newClientMap()
	})
	slog.Debug("Returning existing client map")
	return ClientMap
}

func newClientMap() *Map {
	return &Map{
		Map:   make(map[string]*Client),
		Mutex: sync.RWMutex{},
	}
}

func (m *Map) AddClient(c *Client) {
	m.Map[c.ID] = c
}

func (m *Map) RemoveClient(id string) {
	delete(m.Map, id)
}

func (c *Client) SendMessage(targetID string, messageType int, message string) {
	clientMap := GetClientMap()

	clientMap.Mutex.Lock()
	targetClient, exists := clientMap.Map[targetID]
	clientMap.Mutex.Unlock()

	if !exists {
		slog.Warn("Target client not found", "ID", targetID)
		return
	}

	err := targetClient.Conn.WriteMessage(messageType, []byte(message))
	if err != nil {
		slog.Warn("Error sending message", "ID", targetID, "error", err)
	}
}

// RandomizeTarget sets the random target to the client and returns its ID
func (c *Client) RandomizeTarget() string {
	clientMap := GetClientMap()
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	var id string
	// Get a random target
	for id = range clientMap.Map {
		if id == c.ID {
			continue
		}
		c.ConnectedTo = id
		break
	}
	return id
}
