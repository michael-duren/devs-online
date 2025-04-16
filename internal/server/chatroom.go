package server

import (
	"sync"

	"github.com/charmbracelet/log"
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
	// key: conn, value: participant{username, online}
	clients map[*websocket.Conn]messages.Participant

	// history
	messages []messages.Message

	// events
	register   chan *Client
	unregister chan *websocket.Conn
	broadcast  chan messages.Message

	mutex sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:    make(map[*websocket.Conn]messages.Participant),
		messages:   []messages.Message{},
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
			cr.clients[client.conn] = messages.Participant{
				Username: client.username,
				Online:   true,
			}
			cr.mutex.Unlock()

			initMsg := messages.NewInitMessage(cr.messages, cr.GetParticipants())
			if err := client.conn.WriteJSON(initMsg); err != nil {
				log.Warnf("unable to write to client after joining. error: %v", err)
				cr.mutex.Lock()
				delete(cr.clients, client.conn)
				cr.mutex.Unlock()
				break
			}

			joinMsg := messages.NewJoinMessage(client.username)
			// keep track of msg history
			cr.messages = append(cr.messages, *joinMsg)

			// broadcast to clients
			cr.broadcast <- *joinMsg
		case webConn := <-cr.unregister:
			cr.mutex.Lock()
			p := cr.clients[webConn]
			delete(cr.clients, webConn)
			cr.mutex.Unlock()

			leaveMsg := messages.NewLeaveMessage(
				p.Username,
			)

			cr.broadcast <- *leaveMsg
		case msg := <-cr.broadcast:
			cr.mutex.Lock()

			for c := range cr.clients {
				if err := c.WriteJSON(msg); err != nil {
					c.Close()
					delete(cr.clients, c)
				}
			}
		}
	}
}

func (cr *ChatRoom) GetParticipants() []messages.Participant {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	participants := make([]messages.Participant, 0, len(cr.clients))
	for _, p := range cr.clients {
		participants = append(participants, p)
	}

	return participants
}
