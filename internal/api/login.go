package api

import (
	"net/http"
	"omar-kada/autonas/internal/logger"
	"omar-kada/autonas/internal/storage"
)

// LoginHandler processes login Http requests
type LoginHandler struct {
	log   logger.Logger
	store storage.Storage
}

func newLoginHandler(store storage.Storage, log logger.Logger) *LoginHandler {
	return &LoginHandler{
		log:   log,
		store: store,
	}
}

func (h *LoginHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// In a real application, you would validate the credentials against a database
	if username == "admin" && password == "password" {
		w.Write([]byte("Login successful"))
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}
