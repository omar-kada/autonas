// Package server provides implementations of http and ws handlers
package server

import (
	"context"
	"log/slog"
	"net/http"
	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/storage"
	"strconv"
	"time"
)

// Server will listen to requests on a port
type Server interface {
	ListenAndServe(port int) error
	Serve(port int) error
	Shutdown(ctx context.Context)
}

// HTTPServer is responsible for listening and mapping http requests
type HTTPServer struct {
	store            storage.Storage
	loginHandler     *LoginHandler
	websocketHandler *WebsocketHandler
	statusHandler    *StatusHandler
	server           *http.Server
}

// NewServer creates a new http server
func NewServer(store storage.Storage, manager process.Manager) Server {
	return &HTTPServer{
		store:            store,
		loginHandler:     newLoginHandler(store),
		websocketHandler: newWebsocketHandler(store),
		statusHandler:    newStatusHandler(manager),
	}
}

// ListenAndServe initializes handler routes and serves on the given port
func (s *HTTPServer) ListenAndServe(port int) error {
	http.Handle("/", http.FileServer(frontendFileSystem{fs: http.Dir("./frontend/dist")}))
	http.HandleFunc("/api/login", s.loginHandler.handle)
	http.HandleFunc("/api/status", s.statusHandler.handle)
	http.HandleFunc("/ws", s.websocketHandler.handle)
	slog.Info("Server starting on ", "port", port)

	s.server = &http.Server{
		Addr:              ":" + strconv.Itoa(port),
		ReadHeaderTimeout: 3 * time.Second,
	}

	return s.server.ListenAndServe()
}

// Serve initializes routes from generated api and serves on the given port
func (s *HTTPServer) Serve(port int) error {
	// Create a new serve mux
	mux := http.NewServeMux()

	// Add frontend file server
	mux.Handle("/", http.FileServer(frontendFileSystem{fs: http.Dir("./frontend/dist")}))

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	myHandler := NewHandler(s.store)
	strict := api.NewStrictHandler(myHandler, []api.StrictMiddlewareFunc{})

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(strict, mux)

	s.server = &http.Server{
		Handler:           h,
		Addr:              ":" + strconv.Itoa(port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	slog.Info("Server starting on ", "port", port)

	// And we serve HTTP until the world ends.
	return s.server.ListenAndServe()
}

// Shutdown closes the server
func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
