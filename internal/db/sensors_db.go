package db

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DatabaseManager) insertTopic(topic *models.Sensor) (*models.Sensor, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		result, err := topicsCollection.InsertOne(context.TODO(), topic)
		if err != nil {
			return nil, err
		}

		topicId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for node")
		}
		topic.ID = topicId
	}

	return topic, nil
}

func (d *DatabaseManager) insertTopics(topics []interface{}) []interface{} {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		result, err := topicsCollection.InsertMany(context.TODO(), topics)
		if err != nil {
			slog.Error("fail to insert topics", "topcis", topics)
			return nil
		}
		return result.InsertedIDs
	}
	return nil
}

func (d *DatabaseManager) DeactivateTopics(sensors []*models.Sensor) {
	filter := buildTopicFilter(sensors)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_active", Value: false}}}}
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		_, err := topicsCollection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			slog.Error("fail to update topics", "topcis", sensors)
			return
		}
	}
}

func (d *DatabaseManager) ActivateTopics(sensors []*models.Sensor) {
	filter := buildTopicFilter(sensors)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "is_active", Value: true}}}}
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(TOPICS_COLLECTION)
		_, err := topicsCollection.UpdateMany(context.TODO(), filter, update)
		if err != nil {
			slog.Error("fail to update topics", "topcis", sensors)
			return
		}
	}
}

func (d *DatabaseManager) FindTopics(filter primitive.D) ([]models.Sensor, error) {
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

		if err = cursor.All(context.TODO(), &results); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func (d *DatabaseManager) FindAllTopics() ([]models.Sensor, error) {
	return d.FindTopics(bson.D{{}})
}
