package handlers

import "mimir/internal/models"

type XMLHandler struct{}

func NewXMLHandler() *XMLHandler {
	return &XMLHandler{}
}

func (p *XMLHandler) HandleMessage(msg Message) error {
	panic("Missing implementation")
}

func (p *XMLHandler) SetReadingsChannel(readingsChannel chan models.SensorReading) {
	panic("Missing implementation")
}

func (p *XMLHandler) GetConfigFilename() string {
	panic("Missing implementation")
}

func (p *XMLHandler) GetTopic() string {
	panic("Missing implementation")
}

func (p *XMLHandler) GetType() HandlerType {
	return XML_HANDLER
}

func (p *XMLHandler) UpdateFields(map[string]interface{}) error {
	return nil
}
