package db

import (
	"fmt"
	mimir "mimir/internal/mimir/models"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
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

	InfluxDBClient *influxdb3.Client
	MongoDBClient  *mongo.Client
)

func Run() {
	loadTopology()
	fmt.Println(GroupsData)
	fmt.Println(NodesData)
	go processPoints()
}
