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
	topicChannel := make(chan string)

	go mqtt.StartGateway(client, topics, topicChannel)
	go mimir.Run(topicChannel)
	go API.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mqtt.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
