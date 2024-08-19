package mimir

import (
	"encoding/json"
	"fmt"
	"io"
	"mimir/internal/consts"
	"mimir/internal/mimir"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var (
	Topics  TopicManager
	Manager MQTTManager
)

type MQTTManager struct {
	MQTTClient      mqtt.Client
	readingsChannel chan mimir.SensorReading
}

func NewMQTTManager(mqttClient mqtt.Client, readingsChannel chan mimir.SensorReading) *MQTTManager {
	return &MQTTManager{mqttClient, readingsChannel}
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())

	var payload = string(message.Payload()[:])
	jsonDataReader := strings.NewReader(payload)
	decoder := json.NewDecoder(jsonDataReader)
	var profile map[string]interface{}
	for {
		err := decoder.Decode(&profile)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Before Reading")

	id := profile["sensorId"].(string)
	value := profile["data"]

	sensorReading := mimir.SensorReading{SensorID: id, Value: value, Time: time.Now()}
	Manager.readingsChannel <- sensorReading
	fmt.Println("reading sent")
	// sensorReading := mimir.SensorReading{SensorID: id, Value: value, Time: time.Now()}
	// mimir.Data.StoreReading(sensorReading)
	//TODO: send through channel

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

func StartGateway(client mqtt.Client, topics []string, topicChannel chan string, readingsChannel chan mimir.SensorReading, outgoingMessagesChannel chan string) {
	Topics = *NewTopicManager(client, topicChannel)
	Manager = *NewMQTTManager(client, readingsChannel)

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
			newTopicName := <-topicChannel
			Topics.AddTopic(newTopicName)
		}
	}()

	go func() {
		for {
			outgoingMessage := <-outgoingMessagesChannel
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
