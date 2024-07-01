package logging

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(development bool) {
	if development {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, nil))
	}
	slog.SetDefault(logger)
}

func GetLogger() *slog.Logger {
	return logger
}
