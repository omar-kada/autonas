package server

import (
	"log/slog"
	"net/http"
	"omar-kada/autonas/internal/storage"

	"github.com/gorilla/websocket"
)

// WebsocketHandler processes websocker requests
type WebsocketHandler struct {
	store storage.DeploymentStorage
}

func newWebsocketHandler(store storage.DeploymentStorage) *WebsocketHandler {
	return &WebsocketHandler{
		store: store,
	}
}

func (*WebsocketHandler) handle(w http.ResponseWriter, r *http.Request) {
	slog.Debug("WebSocket connection attempt from ", "addr", r.RemoteAddr, "request", r)
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("Upgrade error", "error", err)
		http.Error(w, "error upgrading "+err.Error(), http.StatusInsufficientStorage)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			slog.Error("Read error", "error", err)
			break
		}

		slog.Info("Received ", "payload", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			slog.Error("Write error", "error", err)
			break
		}
	}
}
