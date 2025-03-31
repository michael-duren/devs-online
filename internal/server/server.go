package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

type Server struct {
	port int
}

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// register route
	mux.HandleFunc("/health", s.healthHandler)
	return mux
}

func NewServer(port int) *http.Server {
	NewServer := &Server{
		port: port,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Infof("creating server that listens on %s\n", server.Addr)

	return server
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
