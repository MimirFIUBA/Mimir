package processors

import "mimir/internal/mimir/models"

type XMLProcessor struct{}

func NewXMLProcessor() *XMLProcessor {
	return &XMLProcessor{}
}

func (p *XMLProcessor) ProcessMessage(topic string, payload []byte) error {
	panic("Missing implementation")
}

func (p *XMLProcessor) SetReadingsChannel(readingsChannel chan models.SensorReading) {
	panic("Missing implementation")
}

func (p *XMLProcessor) GetConfigFilename() string {
	panic("Missing implementation")
}
