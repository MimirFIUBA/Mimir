package mimir

import (
	"fmt"
	"mimir/internal/consts"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	Topics            TopicManager
	Manager           MQTTManager
	MessageProcessors *ProcessorRegistry
)

type MQTTManager struct {
	MQTTClient      mqtt.Client
	readingsChannel chan SensorReading
}

func NewMQTTManager(mqttClient mqtt.Client, readingsChannel chan SensorReading) *MQTTManager {
	return &MQTTManager{mqttClient, readingsChannel}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
	// fmt.Printf("binary: %08b\n", message.Payload())

	processor, exists := MessageProcessors.GetProcessor(message.Topic())
	if exists {
		err := processor.ProcessMessage(message.Topic(), message.Payload())
		if err != nil {
			fmt.Println("Error Process Message") //TODO: log error
		}
	}
}

func GetTopics() []string {
	topicTemp := consts.TopicTemp
	topicPH := consts.TopicPH
	topics := []string{topicTemp, topicPH}
	return topics
}

func StartMqttClient() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(consts.Broker)

	return mqtt.NewClient(opts)
}

func (mp *MimirProcessor) StartGateway(client mqtt.Client, topics []string) {
	Topics = *NewTopicManager(client, mp.TopicChannel)
	Manager = *NewMQTTManager(client, mp.ReadingChannel)
	MessageProcessors = NewProcessorRegistry()

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	for _, topic := range topics {
		if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
	}

	go func() {
		for {
			newTopicName := <-mp.TopicChannel
			Topics.AddTopic(newTopicName)
		}
	}()

	go func() {
		for {
			outgoingMessage := <-mp.OutgoingMessagesChannel
			topic := "alert/ph"
			token := client.Publish(topic, 0, false, outgoingMessage)
			token.Wait()

			fmt.Printf("Published topic %s: %s\n", topic, outgoingMessage)
		}
	}()
}

func CloseConnection(client mqtt.Client, topics []string) {
	for _, topic := range topics {
		client.Unsubscribe(topic)
	}
	client.Disconnect(250)
}
