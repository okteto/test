// pkg/logger/logger.go
package logger

import (
	"log/slog"
	"os"
)

// Logger wraps the slog logger
type Logger struct {
	slog *slog.Logger
}

// LogLevel defines the available log levels
type LogLevel string

const (
	LevelInfo  LogLevel = "info"
	LevelDebug LogLevel = "debug"
	LevelError LogLevel = "error"
)

// NewLogger initializes and returns a Logger based on the log level
func NewLogger(level LogLevel) *Logger {
	var logLevel slog.Level
	switch level {
	case LevelDebug:
		logLevel = slog.LevelDebug
	case LevelError:
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: logLevel})
	slogLogger := slog.New(handler)
	return &Logger{slog: slogLogger}
}

// LogInfo logs an info message
func (l *Logger) LogInfo(msg string, keysAndValues ...interface{}) {
	l.slog.Info(msg, keysAndValues...)
}

// LogDebug logs a debug message
func (l *Logger) LogDebug(msg string, keysAndValues ...interface{}) {
	l.slog.Debug(msg, keysAndValues...)
}

// LogError logs an error message
func (l *Logger) LogError(msg string, err error, keysAndValues ...interface{}) {
	l.slog.Error(msg, append(keysAndValues, "error", err)...)
}

// ParseLogLevel converts a string log level into the LogLevel type
func ParseLogLevel(level string) LogLevel {
	switch level {
	case "debug":
		return LevelDebug
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}
