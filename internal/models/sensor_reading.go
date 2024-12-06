package models

import "time"

type SensorReading struct {
	SensorID string                 `json:"sensorId"`
	Topic    string                 `json:"topic"`
	Value    interface{}            `json:"value"`
	Time     time.Time              `json:"time"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

type SensorValue interface{}
