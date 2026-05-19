package sse

import (
	"encoding/json"
	"sync"
)

type Event struct {
	ID    uint   `json:"id"`
	Type  uint   `json:"type"`
	Title string `json:"title"`
	Body  string `json:"body"`
	RefID *uint  `json:"refID,omitempty"`
}

func (e Event) JSON() string {
	b, _ := json.Marshal(e)
	return string(b)
}

type Hub struct {
	mu      sync.RWMutex
	clients map[uint][]chan Event
}

var Global = &Hub{
	clients: make(map[uint][]chan Event),
}

func (h *Hub) Subscribe(userID uint) chan Event {
	ch := make(chan Event, 10)
	h.mu.Lock()
	h.clients[userID] = append(h.clients[userID], ch)
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(userID uint, ch chan Event) {
	h.mu.Lock()
	defer h.mu.Unlock()
	channels := h.clients[userID]
	for i, c := range channels {
		if c == ch {
			h.clients[userID] = append(channels[:i], channels[i+1:]...)
			break
		}
	}
	if len(h.clients[userID]) == 0 {
		delete(h.clients, userID)
	}
	close(ch)
}

func (h *Hub) Send(userID uint, event Event) {
	h.mu.RLock()
	channels := make([]chan Event, len(h.clients[userID]))
	copy(channels, h.clients[userID])
	h.mu.RUnlock()
	for _, ch := range channels {
		select {
		case ch <- event:
		default:
		}
	}
}
