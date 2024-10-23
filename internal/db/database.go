package db

import (
	"fmt"
	"log"
	influxdb "mimir/db/influxdb"
	"mimir/db/mongodb"
	"mimir/internal/consts"
	"mimir/internal/models"
	"mimir/triggers"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/gookit/ini/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	SensorsData = SensorsManager{
		idCounter: 0,
		sensors:   make([]models.Sensor, 0),
	}
	NodesData = NodesManager{
		idCounter: 0,
		nodes:     make([]models.Node, 0),
	}
	GroupsData = GroupsManager{
		idCounter: 0,
		groups:    make([]models.Group, 0),
	}

	ActiveTriggers       = make([]triggers.Trigger, 0)
	TriggerFilenamesById = make(map[string]string, 0)

	ReadingsDBBuffer = make([]models.SensorReading, 0)

	Database = DatabaseManager{}
)

type DatabaseManager struct {
	InfluxDBClient DatabaseClient
	MongoDBClient  DatabaseClient
}

type DatabaseClient struct {
	client      interface{}
	isConnected bool
}

func Run() {
	loadTopology()
	go processPoints()
}

func (d *DatabaseManager) ConnectToInfluxDB() (*influxdb3.Client, error) {
	godotenv.Load(ini.String(consts.INFLUX_CONFIGURATION_FILE_CONFIG_NAME))
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
	godotenv.Load(ini.String(consts.MONGO_CONFIGURATION_FILE_CONFIG_NAME))
	client, err := mongodb.Connect()
	if err != nil {
		fmt.Println("Failed to connect to mongo: ", err)
		return nil, err
	} else {
		d.AddMongoClient(client)
	}

	return client, nil
}

func (d *DatabaseManager) AddMongoClient(mongoClient *mongo.Client) {
	d.MongoDBClient = DatabaseClient{mongoClient, true}
}

func (d *DatabaseManager) AddInfluxClient(influxDbClient *influxdb3.Client) {
	d.InfluxDBClient = DatabaseClient{influxDbClient, true}
}

func (d *DatabaseManager) getMongoClient() *mongo.Client {
	if d.MongoDBClient.isConnected {
		client, ok := d.MongoDBClient.client.(*mongo.Client)
		if !ok {
			panic("error getting mongo db client")
		}
		return client
	}
	return nil
}

func (d *DatabaseManager) getInfluxDBClient() *influxdb3.Client {
	if d.InfluxDBClient.isConnected {
		client, ok := d.InfluxDBClient.client.(*influxdb3.Client)
		if !ok {
			panic("error getting influx db client")
		}
		return client
	}
	return nil
}
