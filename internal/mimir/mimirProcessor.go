package mimir

import (
	"fmt"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/mimir/processors"
	"mimir/internal/models"
)

type MimirProcessor struct {
	OutgoingMessagesChannel chan string
	ReadingChannel          chan models.SensorReading
	TopicChannel            chan string
	WsChannel               chan string
}

func NewMimirProcessor() *MimirProcessor {
	topicChannel := make(chan string)
	readingsChannel := make(chan models.SensorReading, 50)
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

func (p *MimirProcessor) RegisterSensor(sensor *models.Sensor) {
	sensor.IsActive = true
	p.TopicChannel <- sensor.Topic
}

func (p *MimirProcessor) publishOutgoingMessages() {
	for {
		outgoingMessage := <-p.OutgoingMessagesChannel
		topic := consts.AlertTopic
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

func StartMimir() *MimirProcessor {
	mp := NewMimirProcessor()
	ActionFactory = models.NewActionFactory(mp.OutgoingMessagesChannel, mp.WsChannel)
	TriggerFactory = models.NewTriggerFactory()
	return mp
}
