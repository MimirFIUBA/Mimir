package consts

const (
	INI_CONFIG_FILENAME = "config/config.ini"

	PROCESSORS_DIR_CONFIG_NAME = "processors_dir"
	PROCESSORS_FILE_SUFFIX     = "_processor.json"
	TRIGGERS_DIR_CONFIG_NAME   = "triggers_dir"
	TRIGGERS_FILE_SUFFIX       = "_trigger.json"

	MONGO_CONFIGURATION_FILE_CONFIG_NAME  = "mongodb_configuration_file"
	INFLUX_CONFIGURATION_FILE_CONFIG_NAME = "influxdb_configuration_file"

	MQTT_BROKER_CONFIG_NAME = "mqtt_broker"

	AlertTopic = "mimir/alerts"
)
