package realtime

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{rooms: make(map[string]map[*websocket.Conn]struct{})}
}

func (h *Hub) Join(ticketID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.rooms[ticketID] == nil {
		h.rooms[ticketID] = make(map[*websocket.Conn]struct{})
	}
	h.rooms[ticketID][conn] = struct{}{}
}

func (h *Hub) Leave(ticketID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.rooms[ticketID]; ok {
		delete(room, conn)
		if len(room) == 0 {
			delete(h.rooms, ticketID)
		}
	}
}

func (h *Hub) Broadcast(ticketID string, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for conn := range h.rooms[ticketID] {
		_ = conn.WriteMessage(websocket.TextMessage, data)
	}
}
