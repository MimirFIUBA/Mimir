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

func loadConfigFile() {
	err := ini.LoadExists("config/config.ini")
	if err != nil {
		fmt.Println("Error loading config file, loading default values...")
		err = ini.LoadStrings(`
			processors_file = "config/processors.json"
			triggers_file = "config/triggers.json"
			influxdb_configuration_file = "db/test_influxdb.env"
		`)
		if err != nil {
			panic("Could not load initial configuration")
		}
	}
}

func loadConfiguration(mimirProcessor *mimir.MimirProcessor) {
	config.LoadConfig(ini.String("processors_file"))
	config.LoadConfig(ini.String("triggers_file"))
	config.BuildProcessors(mimirProcessor)
	config.BuildTriggers(mimirProcessor)
}

func connectToInfluxDB() (*influxdb3.Client, error) {
	godotenv.Load(ini.String("influxdb_configuration_file"))
	dbClient, err := influxdb.ConnectToInfluxDB()
	if err != nil {
		log.Fatal("Error connecting to InfluxDB ", err)
		return nil, err
	} else {
		mimirDb.InfluxDBClient = dbClient
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
		mimirDb.MongoDBClient = client
	}

	return client, nil
}

func main() {
	fmt.Println("MiMiR starting")

	loadConfigFile()

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

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

	influxClient, err := connectToInfluxDB()
	if err != nil {
		fmt.Println("error connecting to influx ", err)
	} else {
		defer influxClient.Close()
	}

	loadConfiguration(mimirProcessor)
	mimirDb.Run()

	go mimirProcessor.Run()
	go api.Start(mimirProcessor.WsChannel)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
