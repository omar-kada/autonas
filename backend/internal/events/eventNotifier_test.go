package events

import (
	"context"
	"errors"
	"testing"

	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"

	"github.com/stretchr/testify/mock"
)

// MockSend is a mock implementation of the Send function
type MockSend struct {
	mock.Mock
}

// Send is the mock implementation of the Send function
func (m *MockSend) Send(rawURL string, message string) error {
	args := m.Called(rawURL, message)
	return args.Error(0)
}

// MockConfigStore is a mock implementation of the ConfigStore interface
type MockConfigStore struct {
	mock.Mock
	storage.ConfigStore
}

func (m *MockConfigStore) Get() (models.Config, error) {
	args := m.Called()
	return args.Get(0).(models.Config), args.Error(1)
}

func (m *MockConfigStore) IsEventNotificationEnabled(eventType models.EventType) bool {
	args := m.Called(eventType)
	return args.Bool(0)
}

func TestNotificationEventHandler_HandleEvent(t *testing.T) {
	mockSend := new(MockSend)
	mockConfigStore := new(MockConfigStore)

	handler := NewNotificationEventHandler(mockConfigStore).(*NotificationEventHandler)
	handler.Send = mockSend.Send

	t.Run("should send notification when config is valid and event is enabled", func(t *testing.T) {
		cfg := models.Config{
			Settings: models.Settings{
				NotificationURL:   "http://example.com",
				NotificationTypes: []models.EventType{models.EventMisc},
			},
		}
		event := models.Event{
			Type:       models.EventMisc,
			ObjectID:   1,
			ObjectName: "Test Object",
			Msg:        "Test Message",
		}

		mockConfigStore.On("Get").Return(cfg, nil)
		mockSend.On("Send", cfg.Settings.NotificationURL, event.Type.ToEmoji()+" "+event.Type.ToText()+" - [1] Test Object :\n Test Message").Return(nil)

		handler.HandleEvent(context.Background(), event)

		mockConfigStore.AssertExpectations(t)
		mockSend.AssertExpectations(t)
	})

	t.Run("should not send notification when config is invalid", func(t *testing.T) {
		event := models.Event{
			Type:       models.EventMisc,
			ObjectID:   1,
			ObjectName: "Test Object",
			Msg:        "Test Message",
		}

		mockConfigStore.On("Get").Return(models.Config{}, errors.New("config error"))

		handler.HandleEvent(context.Background(), event)

		mockConfigStore.AssertExpectations(t)
		mockSend.AssertNotCalled(t, "Send")
	})

	t.Run("should not send notification when event is not enabled", func(t *testing.T) {
		cfg := models.Config{
			Settings: models.Settings{
				NotificationURL:   "http://example.com",
				NotificationTypes: []models.EventType{},
			},
		}
		event := models.Event{
			Type:       models.EventMisc,
			ObjectID:   1,
			ObjectName: "Test Object",
			Msg:        "Test Message",
		}

		mockConfigStore.On("Get").Return(cfg, nil)

		handler.HandleEvent(context.Background(), event)

		mockConfigStore.AssertExpectations(t)
		mockSend.AssertNotCalled(t, "Send")
	})
}
