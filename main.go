package main

import (
	"fmt"
	mimir "mimir/internal/mimir"
	mqtt "mimir/internal/mqtt"
	API "mimir/internal/restApi"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	topics := mqtt.GetTopics()
	client := mqtt.StartMqttClient()

	go mimir.Run()
	go mqtt.StartGateway(client, topics)
	go API.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mqtt.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
