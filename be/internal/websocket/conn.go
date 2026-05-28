package websocket

import (
	"errors"
	"sync"

	gorillaws "github.com/gorilla/websocket"
)

// errConnClosed is returned when the peer closed the socket or I/O failed.
// Callers must stop reading — gorilla panics on repeated reads after failure.
var errConnClosed = errors.New("websocket closed")

type Conn struct {
	raw    *gorillaws.Conn
	mu     sync.Mutex
	failed bool
}

func newConn(raw *gorillaws.Conn) *Conn {
	return &Conn{raw: raw}
}

func (c *Conn) markFailed() {
	c.mu.Lock()
	c.failed = true
	c.mu.Unlock()
}

func (c *Conn) isFailed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.failed
}

func (c *Conn) WriteJSON(msg ServerMessage) error {
	c.mu.Lock()
	if c.failed {
		c.mu.Unlock()
		return errConnClosed
	}
	err := c.raw.WriteJSON(msg)
	if err != nil {
		c.failed = true
	}
	c.mu.Unlock()
	return err
}

func (c *Conn) WriteMessage(messageType int, data []byte) error {
	c.mu.Lock()
	if c.failed {
		c.mu.Unlock()
		return errConnClosed
	}
	err := c.raw.WriteMessage(messageType, data)
	if err != nil {
		c.failed = true
	}
	c.mu.Unlock()
	return err
}

func (c *Conn) ReadMessage() (int, []byte, error) {
	if c.isFailed() {
		return 0, nil, errConnClosed
	}
	msgType, data, err := c.raw.ReadMessage()
	if err != nil {
		c.markFailed()
		return 0, nil, errConnClosed
	}
	return msgType, data, nil
}

func (c *Conn) Close() error {
	c.markFailed()
	return c.raw.Close()
}
