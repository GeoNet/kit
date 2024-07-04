package slogger

import (
	"log/slog"
	"os"
)

// return a structured slog.Logger by setting env var LOG_LEVEL
// Usage:
// var logger = slogger.New()
// logger.Info("msg", args1, args2)
// logger.Warn("msg", args1, args2)
// logger.Debug("msg", args1, args2)
// logger.Error("msg", args1, args2)
//
//	args: expected in key:value pairs
func New() *slog.Logger {
	var opts *slog.HandlerOptions
	logLevel := os.Getenv("LOG_LEVEL")
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
	return slog.New(slog.NewTextHandler(os.Stdout, opts))
}
