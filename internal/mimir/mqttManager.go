package mimir

import (
	"fmt"
	"log"
	"log/slog"
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
	Topics          map[string]bool
	ReadingsChannel chan models.SensorReading
	newTopicChannel chan string
}

func NewMQTTManager(mqttClient mqtt.Client, readingsChannel chan models.SensorReading, newTopicChannel chan string) *MQTTManager {
	return &MQTTManager{mqttClient, make(map[string]bool), readingsChannel, newTopicChannel}
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

func (m *MQTTManager) CloseConnection() {

	topics := m.GetSubscribedTopics()
	for _, topic := range topics {
		m.MQTTClient.Unsubscribe(topic)
	}
	m.MQTTClient.Disconnect(250)
}

func (m *MQTTManager) AddTopic(topic string) {
	isSubscribed, ok := m.Topics[topic]
	if !ok || !isSubscribed {
		slog.Info("Subcribing to topic " + topic)
		if token := m.MQTTClient.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
		m.Topics[topic] = true
	}
}

func (m *MQTTManager) GetSubscribedTopics() []string {
	var topics []string
	for topic, isSubscribed := range m.Topics {
		if isSubscribed {
			topics = append(topics, topic)
		}
	}
	return topics
}
