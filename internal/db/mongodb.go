package db

import (
	"context"
	mimir "mimir/internal/mimir/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func loadTopology() {
	mongoDBClient := Database.getMongoClient()
	loadGroups(mongoDBClient)
	loadNodes(mongoDBClient)
}

func loadGroups(mongoDBClient *mongo.Client) {
	if mongoDBClient != nil {
		filter := bson.D{}
		topicsCollection := mongoDBClient.Database(MONGO_DB_MIMIR).Collection(GROUPS_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			panic(err)
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []mimir.Group
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		if len(results) > 0 {
			for _, group := range results {
				GroupsData.AddGroup(&group)
			}
		}
	}
}

func loadNodes(mongoDBClient *mongo.Client) {
	if mongoDBClient != nil {

		filter := bson.D{}
		topicsCollection := mongoDBClient.Database(MONGO_DB_MIMIR).Collection(NODES_COLLECTION)
		cursor, err := topicsCollection.Find(context.TODO(), filter)
		if err != nil {
			panic(err)
		} else {
			defer cursor.Close(context.TODO())
		}

		var results []mimir.Node
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}

		if len(results) > 0 {
			for _, node := range results {
				NodesData.AddNode(&node)
			}
		}
	}
}
