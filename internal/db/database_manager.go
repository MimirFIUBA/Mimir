package db

import (
	"context"
	"fmt"
	"log"
	mimir "mimir/internal/mimir/models"

	"github.com/InfluxCommunity/influxdb3-go/influxdb3"
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

func (d *DatabaseManager) insertGroup(group *mimir.Group) (*mimir.Group, error) {
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

func (d *DatabaseManager) insertNode(node *mimir.Node) (*mimir.Node, error) {
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

func (d *DatabaseManager) insertTopic(topic *mimir.Sensor) (*mimir.Sensor, error) {
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

func (d *DatabaseManager) findTopics(filter primitive.D) ([]mimir.Sensor, error) {
	var results []mimir.Sensor
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			return nil, err
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []mimir.Sensor
		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
	}
	return results, nil
}
