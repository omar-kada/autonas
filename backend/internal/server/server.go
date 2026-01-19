// Package server provides implementations of http and ws handlers
package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/process"
	"omar-kada/autonas/internal/storage"

	"github.com/rs/cors"
)

// Server will listen to requests on a port
type Server interface {
	Serve(port int) error
	Shutdown(ctx context.Context)
}

// HTTPServer is responsible for listening and mapping http requests
type HTTPServer struct {
	store            storage.DeploymentStorage
	configStore      storage.ConfigStore
	processSvc       process.Service
	websocketHandler *WebsocketHandler
	server           *http.Server
}

// NewServer creates a new http server
func NewServer(store storage.DeploymentStorage, configStore storage.ConfigStore, service process.Service) Server {
	return &HTTPServer{
		store:            store,
		configStore:      configStore,
		processSvc:       service,
		websocketHandler: newWebsocketHandler(store),
	}
}

func spaHandler(w http.ResponseWriter, r *http.Request) {
	frontDir, _ := filepath.Abs(filepath.Join("frontend", "dist"))
	path := filepath.Join(frontDir, r.URL.Path)

	absPath, err := filepath.Abs(path)
	if err != nil || !strings.HasPrefix(absPath, frontDir) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	_, err = os.Stat(absPath)
	if os.IsNotExist(err) {
		http.ServeFile(w, r, filepath.Join(frontDir, "index.html"))
		return
	}

	http.FileServer(http.Dir(frontDir)).ServeHTTP(w, r)
}

// Serve initializes routes from generated api and serves on the given port
func (s *HTTPServer) Serve(port int) error {
	// Create a new serve mux
	mux := http.NewServeMux()

	// Add frontend file server
	mux.HandleFunc("/ws", s.websocketHandler.handle)
	mux.HandleFunc("/", spaHandler)

	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
	myHandler := NewHandler(s.store, s.configStore, s.processSvc)
	strict := api.NewStrictHandler(myHandler, []api.StrictMiddlewareFunc{})

	// get an `http.Handler` that we can use
	h := api.HandlerFromMux(strict, mux)
	// Set up the CORS filter
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"localhost:*", "127.0.0.1:*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	// Use the CORS filter as a middleware
	h = c.Handler(h)
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
