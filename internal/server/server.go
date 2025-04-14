package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	logger *log.Logger
	Server *http.Server
	hub    *Hub
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("error upgrading request: %v", err)
		return
	}

	client := &Client{
		hub:  s.hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	// allow collection of memroy referenced by the caller by doing all work in
	// new go routiness
	go client.readPump()
	go client.writePump()
}

func (s *Server) ShutdownSockets() {
	// send signal to stop
	s.hub.Stop <- struct{}{}
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
		s.serveWs(w, r)
	})
	return mux
}

func NewServer(port int, logger *log.Logger) *Server {
	Hub := newHub()
	NewServer := &Server{
		logger: logger,
		hub:    Hub,
	}
	// hub manages the different connections and broadcasting
	go NewServer.hub.Run()

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
