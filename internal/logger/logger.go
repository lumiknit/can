package logger

import (
	"log/slog"
	"os"
)

// New creates a new structured logger with the specified release mode.
func New(releaseMode bool) *slog.Logger {
	var handler slog.Handler

	if releaseMode {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return slog.New(handler)
}
