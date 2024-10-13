package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"
	"os"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
	"github.com/gookit/ini/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	MONGO_DB_MIMIR    = "Mimir"
	GROUPS_COLLECTION = "groups"
	NODES_COLLECTION  = "nodes"
	TOPICS_COLLECTION = "topics"
)

type DatabaseManager struct {
	InfluxDBClient DatabaseClient
	MongoDBClient  DatabaseClient
}

type DatabaseClient struct {
	client      interface{}
	isConnected bool
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

func (d *DatabaseManager) insertGroup(group *models.Group) (*models.Group, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		groupsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		result, err := groupsCollection.InsertOne(context.TODO(), group)
		if err != nil {
			fmt.Println("error inserting group ", err)
			return nil, err
		}

		groupId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for group")
		}
		group.ID = groupId
	}
	return group, nil
}

func (d *DatabaseManager) insertNode(node *models.Node) (*models.Node, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		nodesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(NODES_COLLECTION)
		result, err := nodesCollection.InsertOne(context.TODO(), node)
		if err != nil {
			fmt.Println("error inserting group ", err)
			return nil, err
		}

		nodeId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for group")
		}
		node.ID = nodeId
	}

	return node, nil
}

func (d *DatabaseManager) insertTopic(topic *models.Sensor) (*models.Sensor, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		result, err := topicsCollection.InsertOne(context.TODO(), topic)
		if err != nil {
			fmt.Println("error inserting group ", err)
			return nil, err
		}

		topicId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for group")
		}
		topic.ID = topicId.String() //TODO: see if we need to change to primitive.ObjectId
	}

	return topic, nil
}

func (d *DatabaseManager) insertTopics(topics []interface{}) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		_, err := topicsCollection.InsertMany(context.TODO(), topics)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *DatabaseManager) findTopics(filter primitive.D) ([]models.Sensor, error) {
	var results []models.Sensor
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			return nil, err
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []models.Sensor
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *DatabaseManager) SaveProcessor(processor processors.MessageProcessor) {
	jsonString, err := json.MarshalIndent(processor, "", "    ")
	if err != nil {
		fmt.Println("Error ", err)
	}

	fileName := ini.String("processors_dir") + "/" + processor.GetConfigFilename()

	// f, err := os.OpenFile(fileName, os.O_CREATE, os.ModePerm)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	os.WriteFile(fileName, jsonString, os.ModePerm)
}
