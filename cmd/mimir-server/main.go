package main

import (
	"context"
	"fmt"
	influxdb "mimir/db/influxdb"
	"mimir/db/mongodb"
	"mimir/internal/api"
	"mimir/internal/config"
	mimirDb "mimir/internal/db"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"

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

func connectToInfluxDB() {
	godotenv.Load(ini.String("influxdb_configuration_file"))
	dbClient, err := influxdb.ConnectToInfluxDB()
	if err != nil {
		fmt.Println("error connecting to db")
		fmt.Println(err)
	} else {
		defer dbClient.Close()

		mimirDb.InfluxDBClient = dbClient
	}
}

func connectToMongo() (*mongo.Client, error) {
	godotenv.Load(ini.String("mongodb_configuration_file"))
	client, err := mongodb.Connect()
	if err != nil {
		fmt.Println("Failed to connect to mongo: ", err)
		return nil, err
	}
	// } else {
	// 	defer func() {
	// 		if err = client.Disconnect(context.TODO()); err != nil {
	// 			panic(err)
	// 		}
	// 	}()
	// }
	mimirDb.MongoDBClient = client
	return client, nil
}

func main() {
	fmt.Println("MiMiR starting")

	loadConfigFile()

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

	loadConfiguration(mimirProcessor)

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
	connectToInfluxDB()

	go mimirProcessor.Run()
	mimirDb.Run()
	go api.Start(mimirProcessor.WsChannel)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
