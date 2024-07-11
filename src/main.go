package main

import (
	"fmt"
	mimir "mimir/src/mimir"
	mqtt "mimir/src/mqtt"
	API "mimir/src/restApi"
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
