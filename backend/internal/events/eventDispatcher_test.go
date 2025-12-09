package events

import (
	"context"
	"log/slog"
	"omar-kada/autonas/api"
	"omar-kada/autonas/internal/storage"
	"testing"
)

var deploymentID string

func newStore() storage.EventStorage {
	store := storage.NewMemoryStorage()
	dep, _ := store.InitDeployment("test", "")
	deploymentID = dep.Id
	return store
}

func TestNewDefaultDispatcher(t *testing.T) {
	store := storage.NewMemoryStorage()
	dispatcher := NewDefaultDispatcher(store)

	if dispatcher == nil {
		t.Error("Expected non-nil dispatcher")
	}
}

func TestDispatchLevel(t *testing.T) {
	store := newStore()
	dispatcher := NewDefaultDispatcher(store).(*dispatcher)

	ctx := context.WithValue(context.Background(), ObjectID, deploymentID)
	msg := "test message"
	args := []any{"arg1", "arg2"}

	dispatcher.dispatchLevel(ctx, slog.LevelInfo, msg, args...)

	storedEvents := store.GetEvents(deploymentID)
	if len(storedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storedEvents))
	}

	storedEvent := storedEvents[0]
	if storedEvent.Level != api.EventLevel(slog.LevelInfo.String()) {
		t.Errorf("Expected level Info, got %v", storedEvent.Level)
	}
	if storedEvent.Msg != msg {
		t.Errorf("Expected message %s, got %s", msg, storedEvent.Msg)
	}
}

func TestInfo(t *testing.T) {
	store := newStore()
	dispatcher := NewDefaultDispatcher(store)

	ctx := context.WithValue(context.Background(), ObjectID, deploymentID)
	msg := "info message"
	args := []any{"arg1", "arg2"}

	dispatcher.Info(ctx, msg, args...)
	storedEvents := store.GetEvents(deploymentID)

	if len(storedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storedEvents))
	}

	storedEvent := storedEvents[0]
	if storedEvent.Level != api.EventLevel(slog.LevelInfo.String()) {
		t.Errorf("Expected level Info, got %v", storedEvent.Level)
	}
}

func TestError(t *testing.T) {
	store := newStore()
	dispatcher := NewDefaultDispatcher(store)

	ctx := context.WithValue(context.Background(), ObjectID, deploymentID)
	msg := "error message"
	args := []any{"arg1", "arg2"}

	dispatcher.Error(ctx, msg, args...)
	storedEvents := store.GetEvents(deploymentID)

	if len(storedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storedEvents))
	}

	storedEvent := storedEvents[0]
	if storedEvent.Level != api.EventLevel(slog.LevelError.String()) {
		t.Errorf("Expected level Error, got %v", storedEvent.Level)
	}
}

func TestDebug(t *testing.T) {
	store := newStore()
	dispatcher := NewDefaultDispatcher(store)

	ctx := context.WithValue(context.Background(), ObjectID, deploymentID)
	msg := "debug message"
	args := []any{"arg1", "arg2"}

	dispatcher.Debug(ctx, msg, args...)
	storedEvents := store.GetEvents(deploymentID)

	if len(storedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storedEvents))
	}

	storedEvent := storedEvents[0]
	if storedEvent.Level != api.EventLevel(slog.LevelDebug.String()) {
		t.Errorf("Expected level Debug, got %v", storedEvent.Level)
	}
}

func TestWarn(t *testing.T) {
	store := newStore()
	dispatcher := NewDefaultDispatcher(store)

	ctx := context.WithValue(context.Background(), ObjectID, deploymentID)
	msg := "warn message"
	args := []any{"arg1", "arg2"}

	dispatcher.Warn(ctx, msg, args...)
	storedEvents := store.GetEvents(deploymentID)

	if len(storedEvents) != 1 {
		t.Errorf("Expected 1 event, got %d", len(storedEvents))
	}

	storedEvent := storedEvents[0]
	if storedEvent.Level != api.EventLevel(slog.LevelDebug.String()) {
		t.Errorf("Expected level Debug, got %v", storedEvent.Level)
	}
}
