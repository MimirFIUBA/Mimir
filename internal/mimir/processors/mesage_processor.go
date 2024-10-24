package processors

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
	ProcessMessage(topic string, payload []byte) error
	SetReadingsChannel(readingsChannel chan models.SensorReading)
	GetConfigFilename() string
	GetTopic() string
	GetType() HandlerType
	UpdateFields(map[string]interface{}) error
}

type ProcessorRegistry struct {
	processors map[string]MessageHandler
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{processors: make(map[string]MessageHandler)}
}

func (r *ProcessorRegistry) RegisterHandler(topic string, processor MessageHandler) {
	r.processors[topic] = processor
}

func (r *ProcessorRegistry) RemoveHandler(topic string) {
	delete(r.processors, topic)
}

func (r *ProcessorRegistry) GetHandler(topic string) (MessageHandler, bool) {
	processor, exists := r.processors[topic]
	return processor, exists
}
func (r *ProcessorRegistry) GetHandlers() map[string]MessageHandler {
	return r.processors
}
