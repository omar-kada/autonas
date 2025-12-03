package process

import (
	"context"
	"log/slog"
)

type CaptureHandler struct {
	slog.Handler
	deafultHandler slog.Handler
	fn             func(slog.Record)
}

func NewLoggerWith(fn func(record slog.Record)) *slog.Logger {
	return slog.New(NewCaptureHandler(slog.Default().Handler(), fn))
}

func NewCaptureHandler(handler slog.Handler, fn func(slog.Record)) *CaptureHandler {
	return &CaptureHandler{
		Handler:        handler,
		fn:             fn,
		deafultHandler: handler,
	}
}

func (h *CaptureHandler) Handle(ctx context.Context, r slog.Record) error {
	// Save a copy (slog.Record mutates when read)
	h.fn(r.Clone())
	return h.deafultHandler.Handle(ctx, r)
}
