package processors

import (
	"fmt"
	"mimir/internal/mimir/models"
)

type MessageProcessor interface {
	ProcessMessage(topic string, payload []byte) error
	SetReadingsChannel(readingsChannel chan models.SensorReading)
	GetConfigFilename() string
}

type ProcessorRegistry struct {
	processors map[string]MessageProcessor
}

func NewProcessorRegistry() *ProcessorRegistry {
	return &ProcessorRegistry{processors: make(map[string]MessageProcessor)}
}

func (r *ProcessorRegistry) RegisterProcessor(topic string, processor MessageProcessor) {
	fmt.Println("register processor ", topic)
	r.processors[topic] = processor
}

func (r *ProcessorRegistry) GetProcessor(topic string) (MessageProcessor, bool) {
	processor, exists := r.processors[topic]
	return processor, exists
}
func (r *ProcessorRegistry) GetProcessors() map[string]MessageProcessor {
	return r.processors
}
