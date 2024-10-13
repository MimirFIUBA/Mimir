package mimir

import (
	"fmt"
	"log"
	"mimir/internal/consts"
	"mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gookit/ini/v2"
)

var (
	Manager           MQTTManager
	MessageProcessors *processors.ProcessorRegistry
)

type MQTTManager struct {
	MQTTClient      mqtt.Client
	Topics          map[string]Topic
	ReadingsChannel chan models.SensorReading
	newTopicChannel chan string
}

func NewMQTTManager(mqttClient mqtt.Client, readingsChannel chan models.SensorReading, newTopicChannel chan string) *MQTTManager {
	return &MQTTManager{mqttClient, make(map[string]Topic), readingsChannel, newTopicChannel}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())

	processor, exists := MessageProcessors.GetProcessor(message.Topic())
	if exists {
		go func() {
			err := processor.ProcessMessage(message.Topic(), message.Payload())
			if err != nil {
				log.Fatal("Error processing message: ", err)
				fmt.Println("Error Process Message")
			}
		}()
	}
}

func StartMqttClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(ini.String(consts.MQTT_BROKER_CONFIG_NAME))

	return mqtt.NewClient(opts)
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

func (m *MQTTManager) CloseConnection() {

	topics := m.GetSubscribedTopics()
	for _, topic := range topics {
		m.MQTTClient.Unsubscribe(topic)
	}
	m.MQTTClient.Disconnect(250)
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
