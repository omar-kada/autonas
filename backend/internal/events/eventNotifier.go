// Package events handles logic related to events
package events

import (
	"context"
	"fmt"
	"log/slog"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/containrrr/shoutrrr"
)

// NotificationEventHandler is an event handler that sends notifications
type NotificationEventHandler struct {
	configStore storage.ConfigStore
	Send        func(rawURL string, message string) error
}

// NewNotificationEventHandler creates a new notification event handler
func NewNotificationEventHandler(configStore storage.ConfigStore) EventHandler {
	return &NotificationEventHandler{
		configStore: configStore,
		Send:        shoutrrr.Send,
	}
}

// HandleEvent sends a notification for the event
func (h *NotificationEventHandler) HandleEvent(_ context.Context, event models.Event) {
	cfg, err := h.configStore.Get()
	if err != nil {
		slog.Error("can't retrieve config", "error", err)
		return
	}
	if cfg.Settings.NotificationURL != "" {
		h.sendNotification(cfg, event)
	}
}

func (h *NotificationEventHandler) sendNotification(cfg models.Config, event models.Event) {
	if !cfg.IsEventNotificationEnabled(event.Type) {
		return
	}
	message := event.Type.ToEmoji() + " " + event.Type.ToText()
	if event.ObjectID != 0 {
		message += fmt.Sprintf(" - [%v] %v", event.ObjectID, event.ObjectName)
	}
	if event.Msg != "" {
		message += fmt.Sprintf(" :\n %v", event.Msg)
	}

	err := h.Send(cfg.Settings.NotificationURL, message)
	if err != nil {
		slog.Error("can't send notification", "error", err)
	}
}
