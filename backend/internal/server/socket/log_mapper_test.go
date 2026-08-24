package socket

import (
	"log/slog"
	"testing"
	"time"

	"omar-kada/air-compose/api"

	"github.com/stretchr/testify/assert"
)

func TestMapLogLine(t *testing.T) {
	time := time.Now()
	record := slog.Record{
		Level:   slog.LevelInfo,
		Message: "Test message",
		Time:    time,
	}

	record.AddAttrs(slog.String("key1", "value1"), slog.String("key2", "value2"))

	tests := []struct {
		name     string
		input    slog.Record
		expected api.LogLine
	}{
		{
			name:  "Test with info level",
			input: record,
			expected: api.LogLine{
				Msg:   "Test message",
				Level: api.Level(slog.LevelInfo.String()),
				Time:  time,
				Meta: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
		{
			name: "Test with error level",
			input: slog.Record{
				Level:   slog.LevelError,
				Message: "Error message",
				Time:    time,
			},
			expected: api.LogLine{
				Msg:   "Error message",
				Level: api.Level(slog.LevelError.String()),
				Time:  time,
				Meta:  map[string]string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.EqualValues(t, tt.expected, MapLogLine(tt.input))
		})
	}
}
