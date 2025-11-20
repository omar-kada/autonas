package api

import (
	"encoding/json"
	"net/http"
	"omar-kada/autonas/internal/process"
)

// StatusHandler processes login Http requests
type StatusHandler struct {
	processManager process.Manager
}

func newStatusHandler(manager process.Manager) *StatusHandler {
	return &StatusHandler{
		processManager: manager,
	}
}

func (h *StatusHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	containers, err := h.processManager.GetManagedContainers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(containers)
}
