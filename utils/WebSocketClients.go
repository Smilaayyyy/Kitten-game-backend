package utils

import (
	"github.com/gorilla/websocket"
	"sync"
)

// WebSocketClients manages the list of connected WebSocket clients.
type WebSocketClientsManager struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

// NewWebSocketClientsManager creates a new WebSocketClientsManager.
func NewWebSocketClientsManager() *WebSocketClientsManager {
	return &WebSocketClientsManager{
		clients: make(map[*websocket.Conn]bool),
	}
}

// AddClient adds a new WebSocket client to the manager.
func (manager *WebSocketClientsManager) AddClient(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	manager.clients[conn] = true
}

// RemoveClient removes a WebSocket client from the manager.
func (manager *WebSocketClientsManager) RemoveClient(conn *websocket.Conn) {
	manager.mu.Lock()
	defer manager.mu.Unlock()
	delete(manager.clients, conn)
}

// Broadcast sends a message to all connected WebSocket clients.
func (manager *WebSocketClientsManager) Broadcast(message []byte) {
	manager.mu.Lock()
	defer manager.mu.Unlock()

	for client := range manager.clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.Close()
			delete(manager.clients, client)
		}
	}
}

var WebSocketClients = NewWebSocketClientsManager()
