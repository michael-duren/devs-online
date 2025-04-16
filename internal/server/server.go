package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
)

type Server struct {
	logger   *log.Logger
	Server   *http.Server
	ChatRoom *ChatRoom
	secret   string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func (s *Server) serveWs(w http.ResponseWriter, r *http.Request, username string, chatRoom *ChatRoom) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Errorf("error upgrading request: %v", err)
// 		return
// 	}
// 	client := &Client{username: username, conn: conn}
// 	chatRoom.register <- client
// 	go func() {
// 		defer func() {
// 			chatRoom.unregister <- conn
// 			_ = conn.Close()
// 		}()
// 		for {
// 			var msg messages.Message
// 			log.Infof("Waiting to read message from client")
// 			if err := conn.ReadJSON(&msg); err != nil {
// 				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 					log.Warnf("Unexpected close error: %v", err)
// 				} else {
// 					log.Warnf("Read error: %v", err)
// 				}
// 				break
// 			}
// 			log.Infof("RECEIVED CLIENT MESSAGE: Type=%v, Content=%v, Sender=%v",
// 				msg.Type, msg.Content, msg.Sender)
// 			chatRoom.broadcast <- msg
// 		}
// 	}()
// }

func (s *Server) ShutdownSockets() {
	// set 5 second context for server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}
	log.Print("Server exiting")
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// register route
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		secret := r.URL.Query().Get("secret")

		// TODO: Add logic to put in some secret
		if secret != "secret" {
			w.WriteHeader(401)
			return
		}

		s.serveWs(w, r, username, s.ChatRoom)
	})
	return mux
}

func NewServer(port int, logger *log.Logger, secret string) *Server {
	NewServer := &Server{
		logger:   logger,
		ChatRoom: NewChatRoom(),
		secret:   secret,
	}
	// run the chat room
	go NewServer.ChatRoom.Run()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	NewServer.Server = server
	log.Infof("creating server that listens on %s\n", server.Addr)

	return NewServer
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	resp := `
    {"msg": "hello world"}
    `
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write([]byte(resp)); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}
