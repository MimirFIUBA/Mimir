package consts

import "time"

const (
	MQTT_ALERT_TOPIC          = "mimir/alerts"
	MQTT_SUBSCRIPTION_TIMEOUT = 10 * time.Second
	MQTT_QUIESCE              = 5000
	MQTT_MAX_RETRIES          = 5
	MQTT_QOS                  = 2
)
