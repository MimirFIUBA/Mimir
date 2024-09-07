package main

import (
	"fmt"
	"mimir/internal/consts"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(consts.Broker)
	topics := getTopics()

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker: %s", token.Error()))
	}

	for _, topic := range topics {
		if token := client.Subscribe(topic, 0, onMessageReceived); token.Wait() && token.Error() != nil {
			panic(fmt.Sprintf("Error subscribing to topic: %s", token.Error()))
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	CloseConnection(client, topics)
}

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
}

func getTopics() []string {
	topics := []string{consts.AlertTopic}
	return topics
}

func CloseConnection(client mqtt.Client, topics []string) {
	for _, topic := range topics {
		client.Unsubscribe(topic)
	}
	client.Disconnect(250)
}
