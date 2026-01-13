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
	Serve(port int) error
	Shutdown(ctx context.Context)
}

// HTTPServer is responsible for listening and mapping http requests
type HTTPServer struct {
	store            storage.DeploymentStorage
	processSvc       process.Service
	websocketHandler *WebsocketHandler
	server           *http.Server
}

// NewServer creates a new http server
func NewServer(store storage.DeploymentStorage, service process.Service) Server {
	return &HTTPServer{
		store:            store,
		processSvc:       service,
		websocketHandler: newWebsocketHandler(store),
	}
}

// Serve initializes routes from generated api and serves on the given port
func (s *HTTPServer) Serve(port int) error {
	// Create a new serve mux
	mux := http.NewServeMux()

	// Add frontend file server
	mux.Handle("/", http.FileServer(frontendFileSystem{fs: http.Dir("./frontend/dist")}))
	mux.HandleFunc("/ws", s.websocketHandler.handle)

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	myHandler := NewHandler(s.store, s.processSvc)
	strict := api.NewStrictHandler(myHandler, []api.StrictMiddlewareFunc{})

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(strict, mux)

	s.server = &http.Server{
		Handler:           loggingMiddleware(h),
		Addr:              ":" + strconv.Itoa(port),
		ReadHeaderTimeout: 3 * time.Second,
	}
	slog.Info("server starting", "port", port)

	// And we serve HTTP until the world ends.
	return s.server.ListenAndServe()
}

// Shutdown closes the server
func (s *HTTPServer) Shutdown(ctx context.Context) {
	s.server.Shutdown(ctx)
}
