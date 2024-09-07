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

	mimirProcessor := mimir.NewMimirProcessor()

	mimirProcessor.StartGateway(client, topics)
	go mimirProcessor.Run()
	go API.Start(mimirProcessor.WsChannel)
	go API.StartWebSocket()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
