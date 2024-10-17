package mimir

import (
	"fmt"
	"mimir/internal/db"
	mimir "mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"
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

		go func() {
			processReading(reading)
			db.StoreReading(reading)
		}()
	}
}

func CloseConnection() {
	Manager.CloseConnection()
	setTopicsInactive()
}

func setTopicsInactive() {
	db.SensorsData.SetSensorsToInactive()
	db.Database.DeactivateTopics(db.SensorsData.GetSensors())
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

func (p *MimirProcessor) publishOutgoingMessages() {
	for {
		outgoingMessage := <-p.OutgoingMessagesChannel
		topic := "mimir/alert"
		token := Manager.MQTTClient.Publish(topic, 0, false, outgoingMessage)
		token.Wait()

		fmt.Printf("Published topic %s: %s\n", topic, outgoingMessage)
	}
}

func (p *MimirProcessor) StartGateway() {

	client := StartMqttClient()

	Manager = *NewMQTTManager(client, p.ReadingChannel, p.TopicChannel)
	MessageProcessors = processors.NewProcessorRegistry()

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	go func() {
		for {
			newTopicName := <-p.TopicChannel
			Manager.AddTopic(newTopicName)
		}
	}()

	go p.publishOutgoingMessages()
}
