package logging

import (
	"log/slog"
	"os"
)

func CreateAPILogger() *slog.Logger {
	// TODO - Add options to Handler
	// TODO - Add file to log
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}

func CreateRequestLogger() *slog.Logger {
	// TODO - Add options to Handler
	// TODO - Add file to log
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
