package api

import (
	"net/http"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/storage"

	"github.com/gorilla/websocket"
)

// WebsocketHandler processes websocker requests
type WebsocketHandler struct {
	log   logger.Logger
	store storage.Storage
}

func newWebsocketHandler(store storage.Storage, log logger.Logger) *WebsocketHandler {
	return &WebsocketHandler{
		log:   log,
		store: store,
	}
}

func (h *WebsocketHandler) handle(w http.ResponseWriter, r *http.Request) {
	h.log.Debugf("WebSocket connection attempt from: %s, %v", r.RemoteAddr, r)
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
		h.log.Errorf("Upgrade error: %s", err)
		http.Error(w, "error upgrading "+err.Error(), http.StatusInsufficientStorage)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			h.log.Errorf("Read error: %s", err)
			break
		}

		h.log.Infof("Received: %s", p)

		if err := conn.WriteMessage(messageType, p); err != nil {
			h.log.Errorf("Write error: %s", err)
			break
		}
	}
}
