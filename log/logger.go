package log

import (
	"log/slog"
	"os"
)

// a structured logger encapsulating log/slog
var (
	logger   *slog.Logger
	logLevel = os.Getenv("LOG_LEVEL")
)

// init logger
func init() {
	var opts *slog.HandlerOptions
	switch logLevel {
	case "ERROR":
		opts = &slog.HandlerOptions{Level: slog.LevelError}
	case "WARN":
		opts = &slog.HandlerOptions{Level: slog.LevelWarn}
	case "DEBUG":
		opts = &slog.HandlerOptions{Level: slog.LevelDebug}
	default:
		opts = &slog.HandlerOptions{Level: slog.LevelInfo}
	}
	logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
}

// encapsulating func in log/slog
// args: expected in key:value pairs
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// encapsulating func in log/slog
// args: expected in key:value pairs
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// encapsulating func in log/slog
// args: expected in key:value pairs
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// encapsulating func in log/slog
// args: expected in key:value pairs
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}
