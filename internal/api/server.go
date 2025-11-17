// Package api provides implementations of http and ws handlers
package api

import (
	"net/http"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/storage"
	"strconv"
)

// Server is responsile for listening and mapping http requests
type Server struct {
	log              logger.Logger
	store            storage.Storage
	loginHandler     *LoginHandler
	websocketHandler *WebsocketHandler
}

// NewServer creates a new http server
func NewServer(store storage.Storage, log logger.Logger) *Server {
	return &Server{
		log:              log,
		store:            store,
		loginHandler:     newLoginHandler(store, log),
		websocketHandler: newWebsocketHandler(store, log),
	}
}

// ListenAndServe initializes handler routes and serves on the given port
func (s *Server) ListenAndServe(port int) error {
	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)
	http.HandleFunc("/login", s.loginHandler.handle)
	http.HandleFunc("/ws", s.websocketHandler.handle)
	s.log.Infof("Server starting on : %s", port)
	return http.ListenAndServe(":"+strconv.Itoa(port), nil)
}
