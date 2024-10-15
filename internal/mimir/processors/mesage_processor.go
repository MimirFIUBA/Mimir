package processors

import (
	"mimir/internal/mimir/models"
)

type ProcessorType int

const (
	JSON_PROCESSOR ProcessorType = iota
	BYTES_PROCESSOR
	XML_PROCESSOR
)

type MessageProcessor interface {
	ProcessMessage(topic string, payload []byte) error
	SetReadingsChannel(readingsChannel chan models.SensorReading)
	GetConfigFilename() string
	GetTopic() string
	GetType() ProcessorType
	UpdateFields(map[string]interface{}) error
}

type ProcessorRegistry struct {
	processors map[string]MessageProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{processors: make(map[string]MessageProcessor)}
}

func (r *ProcessorRegistry) RegisterProcessor(topic string, processor MessageProcessor) {
	r.processors[topic] = processor
}

func (r *ProcessorRegistry) RemoveProcessor(topic string) {
	delete(r.processors, topic)
}

func (r *ProcessorRegistry) GetProcessor(topic string) (MessageProcessor, bool) {
	processor, exists := r.processors[topic]
	return processor, exists
}
func (r *ProcessorRegistry) GetProcessors() map[string]MessageProcessor {
	return r.processors
}
