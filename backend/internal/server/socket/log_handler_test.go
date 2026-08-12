package socket

import (
	"context"
	"io"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"omar-kada/air-compose/api"
	"omar-kada/air-compose/internal/logs"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSender struct {
	mock.Mock
	sendLogCalls          atomic.Int32
	sendPreviousLogsCalls atomic.Int32
}

func (m *MockSender) SendPreviousLogs(ctx context.Context, logs api.LogMessages) error {
	args := m.MethodCalled("SendPreviousLogs", ctx, logs)
	m.sendPreviousLogsCalls.Add(1)
	return args.Error(0)
}

func (m *MockSender) SendPreviousLogsCount() int32 {
	return m.sendPreviousLogsCalls.Load()
}

func (m *MockSender) SendLog(ctx context.Context, log api.LogLine) error {
	args := m.MethodCalled("SendLog", ctx, log)
	m.sendLogCalls.Add(1)
	return args.Error(0)
}

func (m *MockSender) SendLogCount() int32 {
	return m.sendLogCalls.Load()
}

func (*MockSender) SendEvent(context.Context, api.EventMessage) error { return nil }
func (*MockSender) SendError(context.Context, api.Error) error        { return nil }

func TestLogHandlerStartLogSendsSingleLogLine(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sender := &MockSender{}
	hub := logs.NewHistoryHub(0)
	handler := NewLogHandler(logger, sender, hub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sender.On("SendLog", mock.Anything, mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "hello"
	})).Return(nil).Once()

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{Value: api.StartLogsMessage{PreviousLines: 2}})

	assert.Eventually(t, func() bool {
		if sender.SendLogCount() > 0 {
			return true
		}
		hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0))
		return false
	}, time.Second, 10*time.Millisecond)

	sender.AssertExpectations(t)
}

func TestLogHandlerCancelLogsUnregistersActiveStream(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sender := &MockSender{}
	hub := logs.NewHistoryHub(0)
	handler := NewLogHandler(logger, sender, hub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sender.On("SendLog", mock.Anything, mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "hello"
	})).Return(nil).Once()

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{})

	assert.Eventually(t, func() bool {
		if sender.SendLogCount() > 0 {
			return true
		}
		hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0))
		return false
	}, time.Second, 10*time.Millisecond)

	handler.HandleMessage(ctx, api.ClientMessageEndLogs{})
	hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "world", 0))

	sender.AssertNumberOfCalls(t, "SendLog", 1)
	sender.AssertExpectations(t)
}

func TestLogHandlerStopSendingLogsPreventsFurtherBroadcasts(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sender := &MockSender{}
	hub := logs.NewHistoryHub(0)
	handler := NewLogHandler(logger, sender, hub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	isHello := mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "hello"
	})
	sender.On("SendLog", mock.Anything, isHello).Return(nil).Once()

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{})
	assert.Eventually(t, func() bool {
		if sender.SendLogCount() > 0 {
			return true
		}
		hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0))
		return false
	}, time.Second, 10*time.Millisecond)

	handler.HandleMessage(ctx, api.ClientMessageEndLogs{})
	hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "goodbye", 0))

	sender.AssertNumberOfCalls(t, "SendLog", 1)
	sender.AssertExpectations(t)
}

func TestLogHandlerCanResubscribeAfterCancel(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sender := &MockSender{}
	hub := logs.NewHistoryHub(0)
	handler := NewLogHandler(logger, sender, hub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sender.On("SendLog", mock.Anything, mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "first" && log.Level == api.LevelERROR
	})).Return(nil).Once()

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{})

	assert.Eventually(t, func() bool {
		if sender.SendLogCount() > 0 {
			return true
		}
		hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelError, "first", 0))
		return false
	}, time.Second, 10*time.Millisecond)

	handler.HandleMessage(ctx, api.ClientMessageEndLogs{})

	sender.On("SendLog", mock.Anything, mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "second" && log.Level == api.LevelDEBUG
	})).Return(nil).Once()

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{})

	assert.Eventually(t, func() bool {
		if sender.SendLogCount() > 1 {
			return true
		}
		hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelDebug, "second", 0))
		return false
	}, time.Second, 10*time.Millisecond)

	sender.AssertExpectations(t)
}

func TestLogHandlerSendsSingleAndMultiLineBroadcasts(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	sender := &MockSender{}
	hub := logs.NewHistoryHub(2)
	handler := NewLogHandler(logger, sender, hub)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	sender.On("SendPreviousLogs", mock.Anything, mock.MatchedBy(func(logs api.LogMessages) bool {
		return len(logs) == 2 && logs[0].Msg == "hello" && logs[1].Msg == "world"
	})).Return(nil).Once()

	hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0))
	hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "world", 0))

	handler.HandleMessage(ctx, api.ClientMessageStartLogs{})

	assert.Eventually(t, func() bool {
		return sender.SendPreviousLogsCount() > 0
	}, time.Second, 10*time.Millisecond)

	sender.On("SendLog", mock.Anything, mock.MatchedBy(func(log api.LogLine) bool {
		return log.Msg == "hello" && log.Level == api.LevelINFO
	})).Return(nil).Once()

	hub.Broadcast(slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0))

	assert.Eventually(t, func() bool {
		return sender.SendLogCount() > 0
	}, time.Second, 10*time.Millisecond)

	sender.AssertExpectations(t)
}
