package main

import (
	"fmt"
	API "mimir/internal/api"
	mimir "mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	topics := mimir.GetTopics()
	client := mimir.StartMqttClient()

	mimir.StartGateway(client, topics)
	go mimir.Run()
	go API.Start()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
