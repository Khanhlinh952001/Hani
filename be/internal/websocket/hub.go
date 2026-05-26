package websocket

import (
	"sync"

	"github.com/google/uuid"
)

// Hub tracks active realtime connections (one entry per websocket).
type Hub struct {
	mu    sync.RWMutex
	conns map[string]*RealtimeSession
}

func NewHub() *Hub {
	return &Hub{conns: make(map[string]*RealtimeSession)}
}

func (h *Hub) Register(s *RealtimeSession) {
	h.mu.Lock()
	h.conns[s.connID] = s
	h.mu.Unlock()
}

func (h *Hub) Unregister(connID string) {
	h.mu.Lock()
	delete(h.conns, connID)
	h.mu.Unlock()
}

func (h *Hub) Count() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.conns)
}

func (h *Hub) ActiveForUser(userID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	n := 0
	for _, s := range h.conns {
		if s.userID == userID {
			n++
		}
	}
	return n
}

// Global hub instance for the API process.
var DefaultHub = NewHub()

func newConnID() string {
	return uuid.New().String()
}
