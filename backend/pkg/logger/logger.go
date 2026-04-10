package logger

import (
	"log/slog"
	"os"
	"strings"

	charmlog "github.com/charmbracelet/log"
)

// NewSlogHandler returns a new slog.Handler using charmbracelet/log.
// It configures the logger based on the provided environment.
func NewSlogHandler(env string) slog.Handler {
	var level charmlog.Level
	var formatter charmlog.Formatter

	switch strings.ToLower(env) {
	case "info":
		level = charmlog.InfoLevel
		formatter = charmlog.JSONFormatter
	default:
		level = charmlog.DebugLevel
		formatter = charmlog.TextFormatter
	}

	opts := charmlog.Options{
		ReportTimestamp: true,
		Level:           level,
		Formatter:       formatter,
	}

	// For development, we might want to report caller
	if formatter == charmlog.TextFormatter {
		opts.ReportCaller = true
	}

	handler := charmlog.NewWithOptions(os.Stderr, opts)

	return handler
}

// Setup logger initializes the default slog logger with the charmbracelet handler.
func SetupSlogDefaultLogger(env string) {
	handler := NewSlogHandler(env)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
