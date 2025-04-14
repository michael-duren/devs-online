package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	Server *http.Server
	Hub    *Hub
}

func serveWs(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// register route
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
	return mux
}

func NewServer(port int) *Server {
	Hub := newHub()
	NewServer := &Server{
		Hub: Hub,
	}
	// hub manages the different connections and broadcasting
	go NewServer.Hub.Run()

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
