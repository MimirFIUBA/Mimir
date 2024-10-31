package handlers

import (
	"mimir/internal/models"
)

type HandlerType int

const (
	JSON_HANDLER HandlerType = iota
	BYTES_HANDLER
	XML_HANDLER
)

type MessageHandler interface {
	HandleMessage(msg Message) error
	SetReadingsChannel(readingsChannel chan models.SensorReading)
	GetConfigFilename() string
	GetTopic() string
	GetType() HandlerType
	UpdateFields(map[string]interface{}) error
}
