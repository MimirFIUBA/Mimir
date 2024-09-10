package mimir

import (
	"fmt"
	"mimir/internal/triggers"
)

var (
	Data = DataManager{}
)

type MimirProcessor struct {
	OutgoingMessagesChannel chan string
	ReadingChannel          chan SensorReading
	TopicChannel            chan string
	WsChannel               chan string
}

func NewMimirProcessor() *MimirProcessor {
	topicChannel := make(chan string)
	readingsChannel := make(chan SensorReading)
	outgoingMessagesChannel := make(chan string)
	webSocketMessageChannel := make(chan string)

	Data.topicChannel = topicChannel
	mp := MimirProcessor{outgoingMessagesChannel, readingsChannel, topicChannel, webSocketMessageChannel}
	return &mp
}

func (mp *MimirProcessor) Run() {
	// setInitialData(mp.OutgoingMessagesChannel, mp.WsChannel)

	for {
		reading := <-mp.ReadingChannel

		processReading(reading)

		Data.StoreReading(reading)
	}
}

func processReading(reading SensorReading) {
	//TODO: ver si necesitamos hacer algo aca
	fmt.Printf("Processing reading: %v \n", reading.Value)
}

func (m *MimirProcessor) NewSendMQTTMessageAction(message string) *triggers.SendMessageThroughChannel {
	return &triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: m.OutgoingMessagesChannel}
}
