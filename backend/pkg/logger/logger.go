package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/charmbracelet/log"
)

// NewSlogHandler returns a new slog.Handler using charmbracelet/log.
// It configures the logger based on the provided environment.
func NewSlogHandler(env string) slog.Handler {
	var level log.Level
	var formatter log.Formatter

	switch strings.ToLower(env) {
	case "prod", "production":
		level = log.InfoLevel
		formatter = log.JSONFormatter
	default:
		level = log.DebugLevel
		formatter = log.TextFormatter
	}

	opts := log.Options{
		ReportTimestamp: true,
		Level:           level,
		Formatter:       formatter,
	}

	// For development, we might want to report caller
	if formatter == log.TextFormatter {
		opts.ReportCaller = true
	}

	handler := log.NewWithOptions(os.Stderr, opts)

	return handler
}

// SetupLogger initializes the default slog logger with the charmbracelet handler.
func SetupLogger(env string) {
	handler := NewSlogHandler(env)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}
