package mappers

import (
	"fmt"

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

// MapToPageInfo maps a slice of models.Event to an api.PageInfo, determining if there are more items
// and providing the end cursor for pagination.
func (EventMapper) MapToPageInfo(events []models.Event, limit int) api.PageInfo {
	return MapToPageInfo(events, limit, func(event models.Event) string {
		return fmt.Sprintf("%d", event.ID)
	})
}
