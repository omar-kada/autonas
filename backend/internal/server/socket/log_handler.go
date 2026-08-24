package socket

import (
	"context"
	"log/slog"
	"sync"

	"omar-kada/air-compose/api"
	"omar-kada/air-compose/internal/logs"
)

// LogHandler handles log-related messages and connection events
type LogHandler struct {
	logger *slog.Logger
	sender MessageSender
	logHub logs.Hub

	mu       sync.Mutex
	stopLogs context.CancelFunc
}

// NewLogHandler creates a new LogHandler instance
func NewLogHandler(logger *slog.Logger, sender MessageSender, logHub logs.Hub) *LogHandler {
	return &LogHandler{
		logger: logger,
		sender: sender,
		logHub: logHub,
	}
}

// OnConnect is called when a new connection is established
func (lh *LogHandler) OnConnect(_ context.Context) {
	lh.logger.Debug("[SOCKET] log handler connected")
}

// HandleMessage is called when a new message is received
func (lh *LogHandler) HandleMessage(ctx context.Context, msg any) {
	switch m := msg.(type) {
	case api.ClientMessageStartLogs:
		lh.handleStartLog(ctx, m.Value)
	case api.ClientMessageEndLogs:
		lh.cancelLogs()
	default:
		// Ignore messages that are not related to logs
	}
}

func (lh *LogHandler) handleStartLog(ctx context.Context, msg api.StartLogsMessage) {
	lh.logger.Debug("[SOCKET] started streaming logs", "previousLines", msg.PreviousLines)

	lh.mu.Lock()
	if lh.stopLogs != nil {
		lh.logger.Debug("[SOCKET] already streaming logs")
		lh.mu.Unlock()
		return
	}
	lh.logger.Debug("[SOCKET] registering new log client")
	logClient := lh.logHub.Register()
	var unregisterOnce sync.Once
	cleanup := func() {
		unregisterOnce.Do(func() {
			lh.logHub.Unregister(logClient)
		})
	}
	lh.stopLogs = cleanup
	lh.mu.Unlock()

	go func() {
		defer cleanup()

		for {
			select {
			case lines, ok := <-logClient.Send:
				if !ok {
					lh.logger.Debug("channel closed from source")
					return
				}
				if len(lines) > 1 {

					var messages []api.LogLine
					for _, line := range lines {
						messages = append(messages, MapLogLine(line))
					}
					if err := lh.sender.SendPreviousLogs(ctx, messages); err != nil {
						lh.logger.Error("error sending previous logs ", "err", err)
						return
					}
				} else if len(lines) == 1 {
					if err := lh.sender.SendLog(ctx, MapLogLine(lines[0])); err != nil {
						lh.logger.Error("error reading logs ", "err", err)
						return
					}
				}
			case <-ctx.Done(): // socket closed
				lh.logger.Debug("socket context cancelled")
				return
			}
		}
	}()
}

func (lh *LogHandler) cancelLogs() {
	lh.mu.Lock()
	defer lh.mu.Unlock()
	if lh.stopLogs != nil {
		lh.logger.Debug("[SOCKET] cancelling active log stream")
		lh.stopLogs()
		lh.stopLogs = nil
	}
}
