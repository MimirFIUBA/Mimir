package configuration

import (
	"encoding/json"
	"fmt"
	"mimir/internal/api/logging"
	"os"
)

type Configuration struct {
	Logging logging.LoggerConfiguration `json:"logging"`
	Server  struct {
		Port int `json:"port"`
	} `json:"server"`
}

func GetConfiguration(name string) (*Configuration, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file: %w", err)
	}
	defer file.Close()

	var config Configuration
	if err = json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode configuration: %w", err)
	}

	return &config, nil
}
