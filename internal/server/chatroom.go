package server

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/michael-duren/tui-chat/messages"
)

type Client struct {
	conn     *websocket.Conn
	username string
}

// ChatRoom
// holds the state of the room
type ChatRoom struct {
	// key: conn, value: username
	clients map[*websocket.Conn]string

	messages   []messages.ChatMessage
	register   chan *Client
	unregister chan *websocket.Conn
	broadcast  chan messages.Message
	mutex      sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*websocket.Conn]string),
		messages:   []messages.ChatMessage{},
		register:   make(chan *Client),
		unregister: make(chan *websocket.Conn),
		broadcast:  make(chan messages.Message),
	}
}

func (cr *ChatRoom) Run() {
	for {
		select {
		case client := <-cr.register:
			cr.mutex.Lock()
			cr.clients[client.conn] = client.username
			cr.mutex.Unlock()
		}
	}
}
