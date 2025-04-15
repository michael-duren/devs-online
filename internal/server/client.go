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
	// time allowed to read the next pong msg from the peer.
	pongWait = 60 * time.Second
	// send pings to peer with this period. must be less than pongwait
	pingPeriod = (pongWait * 9) / 10
	// max message size allowed from peer.
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
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

	// user info
	username string

	// buffered cchannel of outbound messages
	send chan []byte
}

// read pump pumps messages from the websocket conn to the hub.
//
// the app runs readpup in a per-connection go routine. the app
// ensures that there is at most one reader on a connection by executing all
// reads from this go routine
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		_ = c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error {
			_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

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
		message = bytes.TrimSpace(bytes.ReplaceAll(message, newLine, space))
		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.conn.Close()
		if err != nil {
			log.Errorf("error closing client connection")
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				log.Warn("unable to set deadline while attempting to write to peer", err)
			}
			if !ok {
				log.Info("connection with peer has been closed")
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// add queued chat messages to the urrent websocket message
			n := len(c.send)
			for range n {
				_, _ = w.Write(newLine)
				_, _ = w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Warn("there was an error attempting to ping the peer", err)
				return
			}
		}
	}
}
