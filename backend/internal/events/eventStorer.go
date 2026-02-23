// Package events handles logic related to events
package events

import (
	"context"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
)

// StoringEventHandler is an event handler that stores events
type StoringEventHandler struct {
	store storage.EventStorage
}

// NewStoringEventHandler creates a new storing event handler
func NewStoringEventHandler(store storage.EventStorage) EventHandler {
	return &StoringEventHandler{
		store: store,
	}
}

// HandleEvent stores the event
func (h *StoringEventHandler) HandleEvent(_ context.Context, event models.Event) {
	if event.ObjectID != 0 {
		h.store.StoreEvent(event)
	}
}
