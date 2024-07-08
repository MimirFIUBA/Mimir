package main

import (
	"fmt"
	API "mimir/src/restApi"
	mqtt "mimir/src/mqtt"
	mimir "mimir/src/mimir"
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