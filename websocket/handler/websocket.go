package handler

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

type WSHandler struct {
	Upgrader websocket.Upgrader
}

type Client struct {
	ID   string
	Conn *websocket.Conn
}

var (
	clients   = make(map[string]*Client)
	clientsMu sync.Mutex
)

func addClient(id string, conn *websocket.Conn) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	clients[id] = &Client{ID: id, Conn: conn}
}

func removeClient(id string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()
	delete(clients, id)
}

func NewWSHandler() *WSHandler {
	return &WSHandler{
		Upgrader: websocket.Upgrader{},
	}
}

func (h *WSHandler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.websocketHandler(w, r)
	}
}

func (h *WSHandler) websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	id := r.URL.Query().Get("id")
	if id == "" {
		fmt.Println("No ID provided")
		return
	}
	target := r.URL.Query().Get("target")
	if target == "" {
		fmt.Println("No target provided")
		return
	}

	addClient(id, conn)
	defer removeClient(id)

	fmt.Printf("Client %s connected\n", id)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Send the message to the target client
		sendMessage(target, messageType, string(message)+" from "+id)
	}
}

func sendMessage(targetID string, messageType int, message string) {
	clientsMu.Lock()
	defer clientsMu.Unlock()

	targetClient, exists := clients[targetID]
	if !exists {
		fmt.Printf("Target client %s not found\n", targetID)
		return
	}

	err := targetClient.Conn.WriteMessage(messageType, []byte(message))
	if err != nil {
		fmt.Printf("Error sending message to %s: %v\n", targetID, err)
	}
}
