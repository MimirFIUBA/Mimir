package config

import (
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/mimir"

	"github.com/gookit/config"
	"github.com/gookit/ini/v2"
)

const DEFAULT_CONFIGURATION string = `
	processors_dir = "config/processors"
	triggers_dir = "config/triggers"
	influxdb_configuration_file = "db/influxdb/test_influxdb.env"
	mongodb_configuration_file = "db/mongodb/test_mongodb.env"
	mqtt_broker = "tcp://broker.emqx.io:1883"`

func LoadIni() {
	err := ini.LoadExists(consts.INI_CONFIG_FILENAME)
	if err != nil {
		slog.Warn("Error loading config file, loading default values...")
		err = ini.LoadStrings(DEFAULT_CONFIGURATION)
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

func BuildInitialConfiguration(mimirEngine *mimir.MimirEngine) {
	BuildHandlers(mimirEngine)
	BuildTriggers(mimirEngine)
	LoadVariables()
}
