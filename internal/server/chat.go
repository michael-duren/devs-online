package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/michael-duren/tui-chat/messages"
)

// Constants for connection management
const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 4096
)

// Client structure with send channel
type Client struct {
	username string
	conn     *websocket.Conn
	chatRoom *ChatRoom
	send     chan messages.Message // Buffered channel for outbound messages
}

// readPump reads messages from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.chatRoom.unregister <- c.conn
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg messages.Message
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Warnf("Unexpected close error: %v", err)
			} else {
				log.Warnf("Read error: %v", err)
			}
			break
		}

		log.Infof("RECEIVED CLIENT MESSAGE: Type=%v, Content=%v, Sender=%v",
			msg.Type, msg.Content, msg.Sender)
		c.chatRoom.broadcast <- msg
		log.Infof("after recieve client message")
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The chat room closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(message); err != nil {
				log.Errorf("Error writing message: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request, username string, chatRoom *ChatRoom) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("error upgrading request: %v", err)
		return
	}

	client := &Client{
		username: username,
		conn:     conn,
		chatRoom: chatRoom,
		send:     make(chan messages.Message, 256), // Buffered channel
	}

	// Register with chat room
	chatRoom.register <- client

	// Start the pumps in new goroutines
	go client.writePump()
	go client.readPump()
}

type ChatRoom struct {
	// Maps for both client objects and participants
	clients      map[*websocket.Conn]*Client
	participants map[*websocket.Conn]messages.Participant
	messages     []messages.Message
	register     chan *Client
	unregister   chan *websocket.Conn
	broadcast    chan messages.Message
	mutex        sync.Mutex
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		clients:      make(map[*websocket.Conn]*Client),
		participants: make(map[*websocket.Conn]messages.Participant),
		messages:     []messages.Message{},
		register:     make(chan *Client),
		unregister:   make(chan *websocket.Conn),
		broadcast:    make(chan messages.Message),
	}
}

func (cr *ChatRoom) Run() {
	for {
		log.Infof("=== STARTING NEW CHATROOM ITERATION ===")
		log.Infof("=== WAITING FOR SELECT ===")

		select {
		case client := <-cr.register:
			log.Infof("=== HANDLING REGISTER ===")
			log.Infof("attemtping to register user: %v", client.username)
			cr.mutex.Lock()
			log.Infof("=== GOT MUTEX LOCK FOR REGISTER ===")
			cr.clients[client.conn] = client
			cr.participants[client.conn] = messages.Participant{
				Username: client.username,
				Online:   true,
			}
			log.Infof("current clients: %v", cr.clients)
			cr.mutex.Unlock()
			log.Infof("=== RELEASED MUTEX LOCK FOR REGISTER ===")

			initMsg := messages.NewInitMessage(cr.messages, cr.GetParticipants())
			log.Infof("=== SENDING INIT MESSAGE ===")
			if err := client.conn.WriteJSON(initMsg); err != nil {
				log.Warnf("unable to write to client after joining. error: %v", err)
				cr.mutex.Lock()
				log.Infof("=== GOT MUTEX LOCK FOR REGISTER ERROR ===")
				delete(cr.clients, client.conn)
				delete(cr.participants, client.conn)
				cr.mutex.Unlock()
				log.Infof("=== RELEASED MUTEX LOCK FOR REGISTER ERROR ===")
				log.Infof("=== REGISTER ERROR COMPLETE ===")
				break
			}

			log.Infof("wrote init msg to client: %v", initMsg)
			joinMsg := messages.NewJoinMessage(client.username)
			// keep track of msg history
			cr.messages = append(cr.messages, *joinMsg)

			// IMPORTANT: Instead of sending to broadcast channel, handle it directly
			log.Infof("=== DIRECTLY HANDLING JOIN MESSAGE ===")

			// Broadcast join message to all clients
			cr.mutex.Lock()
			clientsList := make([]*websocket.Conn, 0, len(cr.clients))
			for c := range cr.clients {
				clientsList = append(clientsList, c)
			}
			cr.mutex.Unlock()

			for _, c := range clientsList {
				if err := c.WriteJSON(joinMsg); err != nil {
					log.Errorf("Error writing join message to client: %v", err)
					_ = c.Close()
					cr.mutex.Lock()
					delete(cr.clients, c)
					delete(cr.participants, c)
					cr.mutex.Unlock()
				}
			}

			log.Infof("=== JOIN MESSAGE BROADCAST COMPLETE ===")
			log.Infof("=== REGISTER COMPLETE ===")

		case webConn := <-cr.unregister:
			log.Infof("=== HANDLING UNREGISTER ===")
			cr.mutex.Lock()
			log.Infof("=== GOT MUTEX LOCK FOR UNREGISTER ===")
			p, exists := cr.participants[webConn]
			if !exists {
				log.Warnf("=== CONNECTION NOT FOUND IN PARTICIPANTS MAP ===")
				cr.mutex.Unlock()
				log.Infof("=== RELEASED MUTEX LOCK FOR UNREGISTER - NOT FOUND ===")
				continue
			}
			delete(cr.clients, webConn)
			delete(cr.participants, webConn)
			cr.mutex.Unlock()
			log.Infof("=== RELEASED MUTEX LOCK FOR UNREGISTER ===")

			leaveMsg := messages.NewLeaveMessage(
				p.Username,
			)
			log.Infof("=== BROADCASTING LEAVE MESSAGE ===")
			cr.broadcast <- *leaveMsg
			log.Infof("=== UNREGISTER COMPLETE ===")

		case msg := <-cr.broadcast:
			log.Infof("=== HANDLING BROADCAST ===")
			log.Infof("BROADCAST DETAILS: Type=%v, Content=%v, Sender=%v",
				msg.Type, msg.Content, msg.Sender)
			log.Infof("Total messages before broadcast: %d", len(cr.messages))

			cr.mutex.Lock()
			log.Infof("=== GOT MUTEX LOCK FOR BROADCAST ===")

			// Add the message to the chat room's message history
			cr.messages = append(cr.messages, msg)
			log.Infof("=== APPENDED MESSAGE TO HISTORY ===")
			log.Infof("Total messages after broadcast: %d", len(cr.messages))

			// Create a copy of clients for safe iteration
			log.Infof("=== CREATING CLIENT LIST ===")
			clientsList := make([]*websocket.Conn, 0, len(cr.clients))
			for c := range cr.clients {
				clientsList = append(clientsList, c)
			}
			log.Infof("broadcasting msg: %v to %d clients", msg, len(clientsList))
			cr.mutex.Unlock()
			log.Infof("=== RELEASED MUTEX LOCK FOR BROADCAST ===")

			// Now iterate over the copy, to avoid modifying the map during iteration
			for _, c := range clientsList {
				log.Infof("=== SENDING TO CLIENT ===")
				if err := c.WriteJSON(msg); err != nil {
					log.Errorf("Error writing to client: %v", err)
					log.Infof("=== CLOSING FAILED CLIENT CONNECTION ===")
					_ = c.Close()
					cr.mutex.Lock()
					log.Infof("=== GOT MUTEX LOCK FOR CLIENT REMOVAL ===")
					delete(cr.clients, c)
					delete(cr.participants, c)
					cr.mutex.Unlock()
					log.Infof("=== RELEASED MUTEX LOCK FOR CLIENT REMOVAL ===")
				} else {
					log.Infof("=== SUCCESSFULLY SENT MESSAGE TO CLIENT ===")
				}
			}
			log.Infof("=== BROADCAST COMPLETE ===")
		}

		log.Infof("=== COMPLETED ITERATION ===")
	}
}

func (cr *ChatRoom) GetParticipants() []messages.Participant {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	participants := make([]messages.Participant, 0, len(cr.participants))
	for _, p := range cr.participants {
		participants = append(participants, p)
	}
	return participants
}
