package events

import (
	"context"
	"omar-kada/autonas/models"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockEventStorage is a mock implementation of EventStorage

type MockEventStorage struct {
	mock.Mock
}

// StoreEvent is a mock implementation of the StoreEvent method
func (m *MockEventStorage) StoreEvent(event models.Event) error {
	args := m.Called(event)
	return args.Error(0)
}

// GetEvents is a mock implementation of the GetEvents method
func (m *MockEventStorage) GetEvents(objectID uint64) ([]models.Event, error) {
	args := m.Called(objectID)
	return args.Get(0).([]models.Event), args.Error(1)
}

func TestStoringEventHandler_HandleEvent(t *testing.T) {
	// Create a new mock event storage
	mockStorage := new(MockEventStorage)

	// Create a new storing event handler with the mock storage
	handler := NewStoringEventHandler(mockStorage)

	// Test case 1: Event with ObjectID 0 should not be stored
	event1 := models.Event{
		ID:         1,
		Type:       models.EventMisc,
		Msg:        "Test event 1",
		Time:       time.Now(),
		ObjectID:   0,
		ObjectName: "Test object 1",
	}

	// Expect the StoreEvent method to not be called
	mockStorage.On("StoreEvent", event1).Return(nil).Once()

	// Call the HandleEvent method with the event
	handler.HandleEvent(context.Background(), event1)

	// Assert that the StoreEvent method was not called
	mockStorage.AssertNotCalled(t, "StoreEvent", event1)

	// Test case 2: Event with ObjectID 1 should be stored
	event2 := models.Event{
		ID:         2,
		Type:       models.EventMisc,
		Msg:        "Test event 2",
		Time:       time.Now(),
		ObjectID:   1,
		ObjectName: "Test object 2",
	}

	// Expect the StoreEvent method to be called with the event
	mockStorage.On("StoreEvent", event2).Return(nil).Once()

	// Call the HandleEvent method with the event
	handler.HandleEvent(context.Background(), event2)

	// Assert that the StoreEvent method was called with the event
	mockStorage.AssertCalled(t, "StoreEvent", event2)
}
