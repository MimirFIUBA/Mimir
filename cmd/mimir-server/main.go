package main

import (
	"context"
	"fmt"
	"log"
	influxdb "mimir/db/influxdb"
	"mimir/db/mongodb"
	"mimir/internal/api"
	"mimir/internal/config"
	mimirDb "mimir/internal/db"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/gookit/ini/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func connectToInfluxDB() (*influxdb3.Client, error) {
	godotenv.Load(ini.String("influxdb_configuration_file"))
	dbClient, err := influxdb.ConnectToInfluxDB()
	if err != nil {
		log.Fatal("Error connecting to InfluxDB ", err)
		return nil, err
	} else {
		mimirDb.Database.AddInfluxClient(dbClient)
		return dbClient, nil
	}
}

func connectToMongo() (*mongo.Client, error) {
	godotenv.Load(ini.String("mongodb_configuration_file"))
	client, err := mongodb.Connect()
	if err != nil {
		fmt.Println("Failed to connect to mongo: ", err)
		return nil, err
	} else {
		mimirDb.Database.AddMongoClient(client)
	}

	return client, nil
}

func main() {
	fmt.Println("MiMiR starting")

	config.LoadIni()
	fmt.Println("Ini loaded")

	mimirProcessor := mimir.NewMimirProcessor()
	fmt.Println("new mimir processor")
	mimirProcessor.StartGateway()
	fmt.Println("gateway started")

	mongoClient, err := connectToMongo()
	if err != nil {
		fmt.Println("error connecting to mongo ", err)
	} else {
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	}
	fmt.Println("connected to mongo")

	influxClient, err := connectToInfluxDB()
	if err != nil {
		fmt.Println("error connecting to influx ", err)
	} else {
		defer influxClient.Close()
	}
	fmt.Println("connected to influx")

	config.LoadConfiguration(mimirProcessor)
	fmt.Println("configuration loaded")
	mimirDb.Run()
	fmt.Println("DB running")

	go mimirProcessor.Run()
	go api.Start(mimirProcessor)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
