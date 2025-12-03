package process

import (
	"context"
	"log/slog"
)

// CaptureHandler captures slog records and passes them to a callback function
type CaptureHandler struct {
	slog.Handler
	deafultHandler slog.Handler
	fn             func(slog.Record)
}

// NewLoggerWith creates a new logger that captures log records and passes them to the given function
func NewLoggerWith(fn func(record slog.Record)) *slog.Logger {
	return slog.New(newCaptureHandler(slog.Default().Handler(), fn))
}

func newCaptureHandler(handler slog.Handler, fn func(slog.Record)) *CaptureHandler {
	return &CaptureHandler{
		Handler:        handler,
		fn:             fn,
		deafultHandler: handler,
	}
}

// Handle captures the log record and passes it to the callback function
func (h *CaptureHandler) Handle(ctx context.Context, r slog.Record) error {
	// Save a copy (slog.Record mutates when read)
	h.fn(r.Clone())
	return h.deafultHandler.Handle(ctx, r)
}
