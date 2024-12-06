package db

import (
	"context"
	"fmt"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (d *DatabaseManager) InsertMessage(message *models.Message) (*models.Message, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(MESSAGES_COLLECTION)
		result, err := topicsCollection.InsertOne(context.TODO(), message)
		if err != nil {
			return nil, err
		}

		messageId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for message")
		}
		message.Id = messageId
	}

	return message, nil
}

func (d *DatabaseManager) FindAllMessages() ([]models.Message, error) {
	var results []models.Message
	mongoClient := d.getMongoClient()
	filter := bson.D{{}}
	if mongoClient != nil {
		findOptions := options.Find()
		findOptions.SetLimit(20)

		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(MESSAGES_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter, findOptions)
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

func (d *DatabaseManager) GetMessageById(id string) (*models.Message, error) {
	var result *models.Message
	mongoClient := d.getMongoClient()
	filter := bson.D{{}}
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(MESSAGES_COLLECTION)
		err := topicsCollection.FindOne(context.TODO(), filter, nil).Decode(result)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (d *DatabaseManager) UpdateMessage(id string, messageUpdate *models.Message) (*models.Message, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		messagesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(MESSAGES_COLLECTION)
		filter := bson.D{{Key: "_id", Value: objectId}}
		update := bson.D{{Key: "$set", Value: messageUpdate}}
		_, err = messagesCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}
		return messageUpdate, nil
	}
	return nil, nil
}
