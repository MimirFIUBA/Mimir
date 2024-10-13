package config

import (
	"fmt"
	"mimir/internal/consts"
	"mimir/internal/mimir"

	"github.com/gookit/config"
	"github.com/gookit/ini/v2"
)

func LoadIni() {
	err := ini.LoadExists(consts.INI_CONFIG_FILENAME)
	if err != nil {
		fmt.Println("Error loading config file, loading default values...")
		err = ini.LoadStrings(`
				processors_file = "config/processors.json"
				triggers_file = "config/triggers.json"
				influxdb_configuration_file = "db/test_influxdb.env"
			`)
		if err != nil {
			panic("Could not load initial configuration")
		}
	}
}

func LoadConfigurationFile(filename string) {
	err := config.LoadFiles(filename)
	if err != nil {
		panic("Error loading configuration file")
	}
}

func BuildInitialConfiguration(mimirProcessor *mimir.MimirProcessor) {
	LoadConfigurationFile(ini.String("processors_file"))
	LoadConfigurationFile(ini.String("triggers_file"))
	BuildProcessors(mimirProcessor)
	BuildTriggers(mimirProcessor)
}
