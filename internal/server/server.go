package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
	"github.com/michael-duren/tui-chat/messages"
)

type Server struct {
	logger *log.Logger
	Server *http.Server
	hub    *Hub
	secret string
}

func (s *Server) serveWs(w http.ResponseWriter, r *http.Request, username string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("error upgrading request: %v", err)
		return
	}

	client := &Client{
		hub:      s.hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		username: username,
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
		var creds messages.Credentials

		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			w.WriteHeader(400)
			log.Info("incorrect format for creds sent to server ", "error: ", err)
			return
		}

		if creds.Secret != s.secret {
			w.WriteHeader(401)
			return
		}

		s.serveWs(w, r, creds.Username)
	})
	return mux
}

func NewServer(port int, logger *log.Logger, secret string) *Server {
	Hub := newHub()
	NewServer := &Server{
		logger: logger,
		hub:    Hub,
		secret: secret,
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
