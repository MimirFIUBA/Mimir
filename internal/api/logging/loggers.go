package logging

import (
	"log/slog"
)

func CreateAPILogger(config *LoggerConfiguration) (*slog.Logger, error) {
	opts := &slog.HandlerOptions{
		AddSource: config.AddSource(),
		Level:     config.GetLevel(),
	}

	file, err := config.GetFile()
	if err != nil {
		return nil, err
	}
	logger := slog.New(slog.NewTextHandler(file, opts))
	return logger, nil
}
