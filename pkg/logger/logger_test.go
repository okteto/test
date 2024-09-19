// pkg/logger/logger_test.go
package logger

import (
	"bytes"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogger_DebugLevel(t *testing.T) {
	// Arrange: Capture the log output
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	slogLogger := slog.New(handler)
	log := &Logger{slog: slogLogger}

	// Act: Log a debug message
	log.LogDebug("Debug log")

	// Assert: The message should appear
	assert.Contains(t, buf.String(), "Debug log")
}

func TestNewLogger_ErrorLevel(t *testing.T) {
	// Arrange: Capture the log output
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError})
	slogLogger := slog.New(handler)
	log := &Logger{slog: slogLogger}

	// Act: Log an error message
	log.LogError("Error log", nil)

	// Assert: The message should appear
	assert.Contains(t, buf.String(), "Error log")
}

func TestLogInfo(t *testing.T) {
	// Arrange: Capture the log output
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo})
	slogLogger := slog.New(handler)
	log := &Logger{slog: slogLogger}

	// Act: Log an info message
	log.LogInfo("Info log")

	// Assert: The message should appear
	assert.Contains(t, buf.String(), "Info log")
}
