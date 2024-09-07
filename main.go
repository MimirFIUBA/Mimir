package main

import (
	"fmt"
	API "mimir/internal/api"
	mimir "mimir/internal/mimir"

	// mqtt "mimir/internal/mqtt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	topics := mimir.GetTopics()
	client := mimir.StartMqttClient()

	topicsChannel := make(chan string)
	readingsChannel := make(chan mimir.SensorReading)
	outgoingMessagesChannel := make(chan string)
	webSocketMessageChannel := make(chan string)

	mimirProcessor := mimir.NewMimirProcessor(topicsChannel, readingsChannel, outgoingMessagesChannel, webSocketMessageChannel)

	mimir.StartGateway(client, topics, topicsChannel, readingsChannel, outgoingMessagesChannel)
	go mimirProcessor.Run()
	go API.Start(webSocketMessageChannel)
	go API.StartWebSocket()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
