package db

import (
	"fmt"
	"log"
	influxdb "mimir/db/influxdb"
	"mimir/db/mongodb"
	mimir "mimir/internal/mimir/models"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/gookit/ini/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	SensorsData = SensorsManager{
		idCounter: 0,
		sensors:   make([]mimir.Sensor, 0),
	}
	NodesData = NodesManager{
		idCounter: 0,
		nodes:     make([]mimir.Node, 0),
	}
	GroupsData = GroupsManager{
		idCounter: 0,
		groups:    make([]mimir.Group, 0),
	}

	ReadingsDBBuffer = make([]mimir.SensorReading, 0)

	Database = DatabaseManager{}
)

func Run() {
	loadTopology()
	go processPoints()
}

func (d *DatabaseManager) ConnectToInfluxDB() (*influxdb3.Client, error) {
	godotenv.Load(ini.String("influxdb_configuration_file"))
	dbClient, err := influxdb.ConnectToInfluxDB()
	if err != nil {
		log.Fatal("Error connecting to InfluxDB ", err)
		return nil, err
	} else {
		d.AddInfluxClient(dbClient)
		return dbClient, nil
	}
}

func (d *DatabaseManager) ConnectToMongo() (*mongo.Client, error) {
	godotenv.Load(ini.String("mongodb_configuration_file"))
	client, err := mongodb.Connect()
	if err != nil {
		fmt.Println("Failed to connect to mongo: ", err)
		return nil, err
	} else {
		d.AddMongoClient(client)
	}

	return client, nil
}
