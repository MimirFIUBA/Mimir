package main

import (
	"fmt"
	API "mimir/internal/api"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"
)



func main() {

	fmt.Println("MiMiR starting")

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

	go mimirProcessor.Run()
	// TODO: Mover a su propio ejecutable, hay que desacoplar
	go API.Start(mimirProcessor.WsChannel)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
