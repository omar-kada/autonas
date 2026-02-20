// Package events handles logic reltaed to events
package events

import (
	"context"
	"log/slog"
	"omar-kada/autonas/internal/storage"
	"omar-kada/autonas/models"
	"time"
)

// ObjectID represent a contextkey for objectID
const ObjectID models.ContextKey = "OBJECT_ID"

// Dispatcher is responsible of processing deployment events
type Dispatcher interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
	Debug(ctx context.Context, msg string, args ...any)
	Warn(ctx context.Context, msg string, args ...any)
}

type dispatcher struct {
	store storage.EventStorage
}

// NewDefaultDispatcher creates a new event dispatcher
func NewDefaultDispatcher(store storage.EventStorage) Dispatcher {
	return &dispatcher{
		store: store,
	}
}

// NewVoidDispatcher creates a new dispatcher that discards all events
// without storing or logging them.
func NewVoidDispatcher() Dispatcher {
	return &dispatcher{}
}

func (d *dispatcher) dispatchLevel(ctx context.Context, level slog.Level, msg string, args ...any) {
	if d.store == nil {
		return
	}
	slog.Log(ctx, level, msg, args...)
	d.store.StoreEvent(models.Event{
		Level:    level,
		Msg:      msg,
		Time:     time.Now(),
		ObjectID: ctx.Value(ObjectID).(uint64),
	})
}

// Dispatch handles processing events
func (d *dispatcher) Info(ctx context.Context, msg string, args ...any) {
	d.dispatchLevel(ctx, slog.LevelInfo, msg, args...)
}

func (d *dispatcher) Error(ctx context.Context, msg string, args ...any) {
	d.dispatchLevel(ctx, slog.LevelError, msg, args...)
}

func (d *dispatcher) Debug(ctx context.Context, msg string, args ...any) {
	d.dispatchLevel(ctx, slog.LevelDebug, msg, args...)
}

func (d *dispatcher) Warn(ctx context.Context, msg string, args ...any) {
	d.dispatchLevel(ctx, slog.LevelWarn, msg, args...)
}
