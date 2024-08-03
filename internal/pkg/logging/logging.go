package logging

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(development bool, fd *os.File, fdErr *os.File) {
	if fd == nil {
		fd = os.Stdout
	}
	if fdErr == nil {
		fdErr = os.Stderr
	}
	if development {
		logger = slog.New(slog.NewTextHandler(fd, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(fdErr, nil))
	}
	slog.SetDefault(logger)
}

func GetLogger() *slog.Logger {
	return logger
}
