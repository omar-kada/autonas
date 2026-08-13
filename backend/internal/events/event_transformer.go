// Package events handles logic related to events
package events

import (
	"context"

	"omar-kada/air-compose/internal/models"
)

// SourceEventTransformer transforms source events into events with additional metadata.
type SourceEventTransformer struct {
	configStore models.ConfigGetter
}

// NewSourceEventTransformer creates a new EventTransformer with the given config store.
func NewSourceEventTransformer(configStore models.ConfigGetter) *SourceEventTransformer {
	return &SourceEventTransformer{
		configStore: configStore,
	}
}

// HandleEvent transforms SourceEvent into an Event
func (t *SourceEventTransformer) HandleEvent(ctx context.Context, srcEvent models.SourceEvent) models.Event {
	cfg := t.configStore.Get()
	event := models.FromSourceEvent(ctx, srcEvent)
	event.IsNotification = cfg.IsEventNotificationEnabled(event.Type)
	return event
}
