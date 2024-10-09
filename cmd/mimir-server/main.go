package main

import (
	"fmt"
	"mimir/db"
	"mimir/internal/config"
	mimirDb "mimir/internal/db"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"syscall"

	"github.com/gookit/ini/v2"
	"github.com/joho/godotenv"
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

func connectToDB() {
	godotenv.Load(ini.String("influxdb_configuration_file"))
	dbClient, err := db.ConnectToInfluxDB()
	if err != nil {
		fmt.Println("error connecting to db")
		fmt.Println(err)
	} else {
		defer dbClient.Close()
		// health, err := dbClient.Health(context.Background())
		// if (err != nil) && health.Status == domain.HealthCheckStatusPass {
		// 	fmt.Println("connectToInfluxDB() error. database not healthy")
		// }

		mimirDb.DBClient = dbClient
	}
}

func main() {
	fmt.Println("MiMiR starting")

	loadConfigFile()

	mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor.StartGateway()

	loadConfiguration(mimirProcessor)

	// connectToDB()

	go mimirProcessor.Run()
	mimirDb.Run()
	// go api.Start(mimirProcessor.WsChannel)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection()

	fmt.Println("Mimir is out of duty, bye!")
}
