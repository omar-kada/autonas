// Package logger handles logging using zap
package logger

import (
	"go.uber.org/zap"
)

// Logger is an abstraction of the logging operations.
type Logger interface {
	Info(msg string, keyvals ...any)
	Error(msg string, keyvals ...any)
	Debug(msg string, keyvals ...any)
	Warn(msg string, keyvals ...any)

	Infof(template string, args ...any)
	Errorf(template string, args ...any)
	Debugf(template string, args ...any)
	Warnf(template string, args ...any)

	With(keyvals ...any) Logger

	Sync() error
}

type zapSugar struct {
	sugar *zap.SugaredLogger
}

// New returns a Logger.
func New(dev bool) Logger {
	var l *zap.Logger
	if dev {
		l, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
	} else {
		l, _ = zap.NewProduction(zap.AddCallerSkip(1))
	}
	return &zapSugar{sugar: l.Sugar()}
}

// Info logs the message with level Info
func (l *zapSugar) Info(msg string, keyvals ...any) { l.sugar.Infow(msg, keyvals...) }

// Error logs the message with level Error
func (l *zapSugar) Error(msg string, keyvals ...any) { l.sugar.Errorw(msg, keyvals...) }

// Debug logs the message with level Debug
func (l *zapSugar) Debug(msg string, keyvals ...any) { l.sugar.Debugw(msg, keyvals...) }

// Warn logs the message with level Warn
func (l *zapSugar) Warn(msg string, keyvals ...any) { l.sugar.Warnw(msg, keyvals...) }

// Infof logs the message with level Info with format
func (l *zapSugar) Infof(template string, args ...any) { l.sugar.Infof(template, args...) }

// Errorf logs the message with level Error with format
func (l *zapSugar) Errorf(template string, args ...any) { l.sugar.Errorf(template, args...) }

// Debugf logs the message with level Debug with format
func (l *zapSugar) Debugf(template string, args ...any) { l.sugar.Debugf(template, args...) }

// Warnf logs the message with level Warn with format
func (l *zapSugar) Warnf(template string, args ...any) { l.sugar.Warnf(template, args...) }

// With adds a variadic number of fields to the logging context.
func (l *zapSugar) With(keyvals ...any) Logger {
	return &zapSugar{sugar: l.sugar.With(keyvals...)}
}

// Sync flushes any buffered log entries.
func (l *zapSugar) Sync() error { return l.sugar.Sync() }
