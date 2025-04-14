package server

type Hub struct {
	clients map[*Client]bool

	// inbound messages from clients
	broadcast chan []byte

	// register requests from client
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Gracefully end conns
	Stop chan struct{}
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Stop:       make(chan struct{}),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			delete(h.clients, client)
			close(client.send)
		case message := <-h.broadcast:
			for c := range h.clients {
				select {
				case c.send <- message:
				default:
					close(c.send)
					delete(h.clients, c)
				}
			}
		case <-h.Stop:
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}

			return
		}
	}
}
