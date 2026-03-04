package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger initializes the global slog logger with rotation and console output
func InitLogger(logDir string) error {
	// Ensure log directory exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile := filepath.Join(logDir, "iptv-server.log")

	// Configure lumberjack for log rotation
	// MaxSize: 10 MB, MaxBackups: 5, MaxAge: 30 days
	rotator := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, // days
		Compress:   true,
	}

	// Create a multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, rotator)

	// Configure slog handler
	handlerOptions := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Optionally, you can add AddSource: true if you want file and line numbers
	}

	// Using TextHandler for readable logs, can be changed to JSONHandler if needed
	handler := slog.NewTextHandler(multiWriter, handlerOptions)

	// Create logger and set it as the default
	logger := slog.New(handler)
	slog.SetDefault(logger)

	return nil
}

// Fatalf is a helper function to log an error message and exit the program.
// slog doesn't have a built-in Fatal method.
func Fatalf(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
