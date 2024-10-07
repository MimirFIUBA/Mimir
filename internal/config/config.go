package config

import (
	"github.com/gookit/config"
)

func LoadConfig(filename string) {
	err := config.LoadFiles(filename)
	if err != nil {
		panic("Error loading configuration file")
	}
}
