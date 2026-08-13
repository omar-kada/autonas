package logs

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTapHandlerHandleInvokesCallbackAndNext(t *testing.T) {
	var output bytes.Buffer
	next := slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})
	var callbackRecords []slog.Record

	handler := NewTapHandler(next, func(_ context.Context, r slog.Record) {
		callbackRecords = append(callbackRecords, r)
	})

	record := slog.NewRecord(time.Unix(1, 0), slog.LevelInfo, "hello world", 0)
	record.AddAttrs(slog.String("key", "value"))

	err := handler.Handle(context.Background(), record)
	require.NoError(t, err)

	require.Len(t, callbackRecords, 1)

	callbackRecord := callbackRecords[0]
	assert.Equal(t, "hello world", callbackRecord.Message)
	assert.Equal(t, slog.LevelInfo, callbackRecord.Level)

	var callbackAttrs []slog.Attr
	callbackRecord.Attrs(func(attr slog.Attr) bool {
		callbackAttrs = append(callbackAttrs, attr)
		return true
	})
	require.Len(t, callbackAttrs, 1)
	assert.Equal(t, "key", callbackAttrs[0].Key)
	assert.Equal(t, "value", callbackAttrs[0].Value.String())

	assert.Contains(t, output.String(), "hello world")
	assert.Contains(t, output.String(), "key=value")
}

func TestTapHandlerDelegatesEnabledAndDecoratorMethods(t *testing.T) {
	var output bytes.Buffer
	next := slog.NewTextHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug})
	handler := NewTapHandler(next, func(context.Context, slog.Record) {})

	assert.True(t, handler.Enabled(context.Background(), slog.LevelInfo))

	wrapped := handler.WithAttrs([]slog.Attr{slog.String("key", "value")})
	wrapped = wrapped.WithGroup("service")

	record := slog.NewRecord(time.Unix(1, 0), slog.LevelInfo, "hello world", 0)
	require.NoError(t, wrapped.Handle(context.Background(), record))

	assert.NotNil(t, wrapped)
	assert.Contains(t, output.String(), "hello world")
	assert.Contains(t, output.String(), "key=value")
}
