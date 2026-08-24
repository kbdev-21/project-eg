package router

import (
	"sync"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/jackc/pgx/v5/pgtype"
)

type WsConnHub struct {
	mu       sync.Mutex
	connsMap map[pgtype.UUID]*websocket.Conn
}

func NewWsConnHub() *WsConnHub {
	return &WsConnHub{
		mu:       sync.Mutex{},
		connsMap: map[pgtype.UUID]*websocket.Conn{},
	}
}

// Register attaches c to userId, closing any previous connection for that user.
func (h *WsConnHub) Register(userId pgtype.UUID, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if old, ok := h.connsMap[userId]; ok && old != nil {
		old.Close()
	}
	h.connsMap[userId] = c
}

// Unregister removes userId's connection only if it still points to c,
// avoiding clobbering a newer connection after a reconnect.
func (h *WsConnHub) Unregister(userId pgtype.UUID, c *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if cur, ok := h.connsMap[userId]; ok && cur == c {
		delete(h.connsMap, userId)
	}
}

// Send writes msg as JSON to userId's connection if one exists.
// The lock is held during WriteJSON so no two goroutines write the same conn.
func (h *WsConnHub) Send(userId pgtype.UUID, msg any) {
	h.mu.Lock()
	defer h.mu.Unlock()

	c, ok := h.connsMap[userId]
	if !ok || c == nil {
		return
	}
	c.WriteJSON(msg)
}
