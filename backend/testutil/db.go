package testutil

import (
	"log/slog"

	"omar-kada/autonas/internal/storage"
)

// NewMemoryStorage instanciates a new memory storage
func NewMemoryStorage() storage.Storage {
	store, err := storage.NewGormStorage(":memory:", 0o000)
	if err != nil {
		slog.Error("couldn't init memory store")
	}
	return store
}
