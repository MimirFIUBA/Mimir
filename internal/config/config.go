package config

import (
	"mimir/internal/mimir"

	"github.com/gookit/config"
	"github.com/gookit/ini/v2"
)

func LoadConfigurationFile(filename string) {
	err := config.LoadFiles(filename)
	if err != nil {
		panic("Error loading configuration file")
	}
}

func LoadConfiguration(mimirProcessor *mimir.MimirProcessor) {
	LoadConfigurationFile(ini.String("processors_file"))
	LoadConfigurationFile(ini.String("triggers_file"))
	BuildProcessors(mimirProcessor)
	BuildTriggers(mimirProcessor)
}
