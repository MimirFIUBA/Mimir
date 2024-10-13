package main

import (
	"context"
	"fmt"
	"mimir/internal/api"
	"mimir/internal/config"
	mimirDb "mimir/internal/db"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("MiMiR starting")

	config.LoadIni()

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

	mongoClient, err := mimirDb.Database.ConnectToMongo()
	if err != nil {
		fmt.Println("error connecting to mongo ", err)
	} else {
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	}

	influxClient, err := mimirDb.Database.ConnectToInfluxDB()
	if err != nil {
		fmt.Println("error connecting to influx ", err)
	} else {
		defer influxClient.Close()
	}

	config.LoadConfiguration(mimirProcessor)
	mimirDb.Run()

	go mimirProcessor.Run()
	go api.Start(mimirProcessor)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
