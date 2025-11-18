// Package api provides implementations of http and ws handlers
package api

import (
	"context"
	"net/http"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/storage"
	"strconv"
	"time"
)

// Server will listen to requests on a port
type Server interface {
	ListenAndServe(port int) error
	Shutdown(ctx context.Context)
}

// HTTPServer is responsible for listening and mapping http requests
type HTTPServer struct {
	log              logger.Logger
	store            storage.Storage
	loginHandler     *LoginHandler
	websocketHandler *WebsocketHandler

	server *http.Server
}

// NewServer creates a new http server
func NewServer(store storage.Storage, log logger.Logger) Server {
	return &HTTPServer{
		log:              log,
		store:            store,
		loginHandler:     newLoginHandler(store, log),
		websocketHandler: newWebsocketHandler(store, log),
	}
}

// ListenAndServe initializes handler routes and serves on the given port
func (s *HTTPServer) ListenAndServe(port int) error {
	http.Handle("/", http.FileServer(frontendFileSystem{fs: http.Dir("./frontend")}))
	http.HandleFunc("/login", s.loginHandler.handle)
	http.HandleFunc("/ws", s.websocketHandler.handle)
	s.log.Infof("Server starting on : %s", port)

	s.server = &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Shutdown closes the server
func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
