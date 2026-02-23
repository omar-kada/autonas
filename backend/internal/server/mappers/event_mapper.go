package mappers

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// EventMapper maps models.Event to api.Event.
type EventMapper struct{}

// Map converts a models.Event to an api.Event.
func (EventMapper) Map(event models.Event) api.Event {
	return api.Event{
		ID:   event.ID,
		Time: event.Time,
		Msg:  event.Msg,
		Type: api.EventType(event.Type),
	}
}
