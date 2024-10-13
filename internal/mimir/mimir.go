package mimir

import (
	"fmt"
	"mimir/internal/db"
	mimir "mimir/internal/mimir/models"
	"mimir/triggers"
)

type MimirProcessor struct {
	OutgoingMessagesChannel chan string
	ReadingChannel          chan mimir.SensorReading
	TopicChannel            chan string
	WsChannel               chan string
}

func NewMimirProcessor() *MimirProcessor {
	topicChannel := make(chan string)
	readingsChannel := make(chan mimir.SensorReading, 50)
	outgoingMessagesChannel := make(chan string)
	webSocketMessageChannel := make(chan string)

	mp := MimirProcessor{
		outgoingMessagesChannel,
		readingsChannel,
		topicChannel,
		webSocketMessageChannel}
	return &mp
}

func (p *MimirProcessor) Run() {
	for {
		reading := <-p.ReadingChannel
		fmt.Println("Reading From channnel in mp ", reading)

		go func() {
			processReading(reading)
			db.StoreReading(reading)
		}()
	}
}

func CloseConnection() {
	Manager.CloseConnection()
}

func processReading(reading mimir.SensorReading) {
	//TODO: ver si necesitamos hacer algo aca
	fmt.Printf("Processing reading: %v \n", reading.Value)
}

// Action creation for simple use
func (p *MimirProcessor) NewSendMQTTMessageAction(message string) triggers.SendMessageThroughChannel {
	return triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: p.OutgoingMessagesChannel}
}

func (p *MimirProcessor) NewSendWebSocketMessageAction(message string) triggers.SendMessageThroughChannel {
	return triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: p.WsChannel}
}

func (p *MimirProcessor) RegisterSensor(sensor *mimir.Sensor) {
	sensor.IsActive = true
	p.TopicChannel <- sensor.Topic
}
