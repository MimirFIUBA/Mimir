package main

import (
	"fmt"
	API "mimir/internal/api"
	mimir "mimir/internal/mimir"
	mqtt "mimir/internal/mqtt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	topics := mqtt.GetTopics()
	client := mqtt.StartMqttClient()

	topicsChannel := make(chan string)
	readingsChannel := make(chan mimir.SensorReading)
	outgoingMessagesChannel := make(chan string)

	mimirProcessor := mimir.NewMimirProcessor(topicsChannel, readingsChannel, outgoingMessagesChannel)

	go mimirProcessor.Run()
	go mqtt.StartGateway(client, topics, topicsChannel, readingsChannel, outgoingMessagesChannel)
	// go mimir.Run(topicChannel)
	go API.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mqtt.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
