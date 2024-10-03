package logging

import (
	"io"
	"log/slog"
	"os"
)

type LoggerConfiguration struct {
	Filename string `json:"file"`
	Level    string `json:"level"`
}

func (lc *LoggerConfiguration) GetFile() (io.Writer, error) {
	fd, err := os.OpenFile(lc.Filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return fd, nil
}

func (lc *LoggerConfiguration) GetLevel() slog.Level {
	switch lc.Level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (lc *LoggerConfiguration) AddSource() bool {
	if lc.GetLevel() == slog.LevelDebug {
		return true
	} else {
		return false
	}
}
