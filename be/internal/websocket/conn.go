package websocket

import (
	"sync"

	gorillaws "github.com/gorilla/websocket"
)

type Conn struct {
	raw *gorillaws.Conn
	mu  sync.Mutex
}

func newConn(raw *gorillaws.Conn) *Conn {
	return &Conn{raw: raw}
}

func (c *Conn) WriteJSON(msg ServerMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.raw.WriteJSON(msg)
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.raw.WriteMessage(messageType, data)
}

func (c *Conn) ReadMessage() (int, []byte, error) {
	return c.raw.ReadMessage()
}

func (c *Conn) Close() error {
	return c.raw.Close()
}
