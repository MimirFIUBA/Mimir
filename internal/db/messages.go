package db

import (
	"context"
	"fmt"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		message.Id = messageId.Hex()
	}

	return message, nil
}
