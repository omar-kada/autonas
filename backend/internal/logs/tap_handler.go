package logs

import (
	"context"
	"log/slog"
)

// TapHandler is a slog.Handler that taps into log records before passing them to the next handler.
type TapHandler struct {
	next  slog.Handler
	onLog func(ctx context.Context, r slog.Record)
}

// NewTapHandler creates a new TapHandler that wraps the given handler and calls the onLog function
// for each log record before passing it to the next handler.
func NewTapHandler(next slog.Handler, onLog func(context.Context, slog.Record)) *TapHandler {
	return &TapHandler{next: next, onLog: onLog}
}

// Enabled reports whether the handler handles records at the given level.
func (h *TapHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

// Handle processes the log record by calling the onLog function and then passing it to the next handler.
func (h *TapHandler) Handle(ctx context.Context, r slog.Record) error {
	h.onLog(ctx, r.Clone())
	return h.next.Handle(ctx, r)
}

// WithAttrs returns a new TapHandler with the given attributes added to the next handler.
func (h *TapHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &TapHandler{next: h.next.WithAttrs(attrs), onLog: h.onLog}
}

// WithGroup returns a new TapHandler with the given group name added to the next handler.
func (h *TapHandler) WithGroup(name string) slog.Handler {
	return &TapHandler{next: h.next.WithGroup(name), onLog: h.onLog}
}
