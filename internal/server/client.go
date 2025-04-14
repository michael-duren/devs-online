package server

import (
	"bytes"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

const (
	// time allowed to write a msg to the peer
	writeWait = 10 * time.Second
	// time allowed to read the next pong msgt from the peer.
	pongWait = 60 * time.Second
	// send pings to peer with this period. must be less than pongwait
	pingPeriod = (pongWait * 9) / 10
	// max message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hub *Hub

	// websocket conn
	conn *websocket.Conn

	// buffered cchannel of outbound messages
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			c.conn.SetReadDeadline(time.Now().Add(pongWait))

			return nil
		})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}
