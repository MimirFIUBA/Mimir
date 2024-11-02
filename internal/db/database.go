package db

import (
	"context"
	"log/slog"
	influxdb "mimir/db/influxdb"
	"mimir/db/mongodb"
	"mimir/internal/consts"
	"mimir/internal/models"
	"mimir/triggers"
	"os"
	"sync"

	"github.com/gookit/ini/v2"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
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

	ReadingsBuffer = ReadingsSyncBuffer{
		make([]models.SensorReading, 0),
		sync.Mutex{},
	}

	Database = DatabaseManager{}
)

type ReadingsSyncBuffer struct {
	buffer []models.SensorReading
	mutex  sync.Mutex
}

func (b *ReadingsSyncBuffer) AddReading(reading models.SensorReading) {
	b.mutex.Lock()
	b.buffer = append(b.buffer, reading)
	b.mutex.Unlock()
}

func (b *ReadingsSyncBuffer) Dump() []models.SensorReading {
	b.mutex.Lock()
	dumpBuffer := make([]models.SensorReading, len(b.buffer))
	copy(dumpBuffer, b.buffer)
	b.buffer = b.buffer[:0]
	b.mutex.Unlock()
	return dumpBuffer
}

type DatabaseManager struct {
	InfluxDBClient DatabaseClient
	MongoDBClient  DatabaseClient
}

type DatabaseClient struct {
	client      interface{}
	isConnected bool
}

func Run(ctx context.Context, wg *sync.WaitGroup) {
	go processPoints(ctx, wg)
}

func (d *DatabaseManager) ConnectToInfluxDB() (influxdb2.Client, error) {
	godotenv.Load(ini.String(consts.INFLUX_CONFIGURATION_FILE_CONFIG_NAME))
	dbClient, err := influxdb.ConnectToInfluxDB()
	if err != nil {
		slog.Error("Error connecting to InfluxDB ", "error", err)
		return nil, err
	}
	d.AddInfluxClient(dbClient)
	return dbClient, nil
}

func (d *DatabaseManager) ConnectToMongo() (*mongo.Client, error) {
	godotenv.Load(ini.String(consts.MONGO_CONFIGURATION_FILE_CONFIG_NAME))
	client, err := mongodb.Connect()
	if err != nil {
		slog.Error("Error connecting to Mongo ", "error", err)
		return nil, err
	}
	d.AddMongoClient(client)
	return client, nil
}

func (d *DatabaseManager) AddMongoClient(mongoClient *mongo.Client) {
	d.MongoDBClient = DatabaseClient{mongoClient, true}
}

func (d *DatabaseManager) AddInfluxClient(influxDbClient influxdb2.Client) {
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

func (d *DatabaseManager) getInfluxDBClient() influxdb2.Client {
	if d.InfluxDBClient.isConnected {
		client, ok := d.InfluxDBClient.client.(influxdb2.Client)
		if !ok {
			panic("error getting influx db client")
		}
		return client
	}
	return nil
}

func (d *DatabaseManager) getInfluxWriteApi() api.WriteAPIBlocking {
	client := d.getInfluxDBClient()
	if client != nil {
		return client.WriteAPIBlocking(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))
	}
	return nil
}

func GetTriggerByName(name string) triggers.Trigger {
	for _, trigger := range ActiveTriggers {
		if trigger.GetName() == name {
			return trigger
		}
	}
	return nil
}
