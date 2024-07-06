package main

import (
	"fmt"
	"mimir/src/consts"
	"os"
	"os/signal"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func onMessageReceived(client mqtt.Client, message mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", message.Payload(), message.Topic())
}

func main() {
	// mqtt topic
	topicTemp := consts.TopicTemp
	topicPH := consts.TopicPH

	opts := mqtt.NewClientOptions()
	opts.AddBroker(consts.Broker)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error connecting to MQTT broker:", token.Error()))
	}

	if token := client.Subscribe(topicTemp, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:", token.Error()))
	}

	fmt.Println("Subscribed to topic:", topicTemp)
	
	if token := client.Subscribe(topicPH, 0, onMessageReceived); token.Wait() && token.Error() != nil {
		panic(fmt.Sprintf("Error subscribing to topic:", token.Error()))
	}

	fmt.Println("Subscribed to topic:", topicPH)

	// Wait for a signal to exit the program gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	client.Unsubscribe(topicTemp)
	client.Unsubscribe(topicPH)
	client.Disconnect(250)
}