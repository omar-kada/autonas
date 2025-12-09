package events

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLoggerWith(t *testing.T) {
	fn := func(_ slog.Record) {}
	logger := NewLoggerWith(fn)
	assert.NotNil(t, logger)
}

func TestHandle(t *testing.T) {
	var capturedRecord slog.Record
	fn := func(record slog.Record) {
		capturedRecord = record
	}
	handler := newCaptureHandler(slog.Default().Handler(), fn)
	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: "test message",
	}
	err := handler.Handle(context.Background(), record)
	assert.NoError(t, err)

	assert.Equal(t, capturedRecord.Message, record.Message)
}
