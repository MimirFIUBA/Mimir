package db

import (
	"context"
	mimir "mimir/internal/mimir/models"

	"go.mongodb.org/mongo-driver/bson"
)

func loadTopology() {
	loadGroups()
	loadNodes()
}

func loadGroups() {
	filter := bson.D{}
	topicsCollection := MongoDBClient.Database("Mimir").Collection("groups")
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

func loadNodes() {
	filter := bson.D{}
	topicsCollection := MongoDBClient.Database("Mimir").Collection("nodes")
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
