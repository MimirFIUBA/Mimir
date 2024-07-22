package mqtt

import (
	"encoding/json"
	"fmt"
	"io"
	"mimir/internal/consts"
	mimir "mimir/internal/mimir"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

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

	id := int(profile["sensorId"].(float64))
	var value mimir.SensorValue
	value = profile["data"]

	sensorReading := mimir.SensorReading{SensorID: id, Value: value, Time: time.Now()}
	mimir.StoreReading(sensorReading)

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

func StartGateway(client mqtt.Client, topics []string) {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	for _, topic := range topics {
		if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
	}
}

func CloseConnection(client mqtt.Client, topics []string) {
	for _, topic := range topics {
		client.Unsubscribe(topic)
	}
	client.Disconnect(250)
}
