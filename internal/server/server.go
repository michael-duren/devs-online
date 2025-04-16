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
		if secret != s.secret {
			w.WriteHeader(401)
			return
		}

		s.serveWs(w, r, username, s.ChatRoom)
	})
	return mux
}

func NewServer(port int, logger *log.Logger, secret string, addr string) *Server {
	NewServer := &Server{
		logger:   logger,
		ChatRoom: NewChatRoom(),
		secret:   secret,
	}

	go NewServer.ChatRoom.Run()
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", addr, port),
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
