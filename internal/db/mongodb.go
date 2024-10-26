package db

import (
	"context"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func LoadTopology() {
	mongoDBClient := Database.getMongoClient()
	loadGroups(mongoDBClient)
	loadNodes(mongoDBClient)
}

func loadGroups(mongoDBClient *mongo.Client) error {
	if mongoDBClient != nil {
		filter := bson.D{}
		topicsCollection := mongoDBClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			return err
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []models.Group
		if err = cursor.All(context.TODO(), &results); err != nil {
			return err
		}

		if len(results) > 0 {
			for _, group := range results {
				GroupsData.AddGroup(&group)
			}
		}
		return nil
	}
	return nil
}

func loadNodes(mongoDBClient *mongo.Client) error {
	if mongoDBClient != nil {
		filter := bson.D{}
		topicsCollection := mongoDBClient.Database(MONGO_DB_MIMIR).Collection(NODES_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			return err
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []models.Node
		if err = cursor.All(context.TODO(), &results); err != nil {
			return err
		}

		if len(results) > 0 {
			for _, node := range results {
				NodesData.AddNode(&node)
			}
		}
	}
	return nil
}
