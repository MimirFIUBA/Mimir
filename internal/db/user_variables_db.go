package db

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *DatabaseManager) findUserVariables(filter primitive.D) ([]UserVariable, error) {
	var results []UserVariable
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		variablesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(VARIABLES_COLLECTION)
		cursor, err := variablesCollection.Find(context.TODO(), filter)
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

func (d *DatabaseManager) insertVariables(variables []interface{}) (*mongo.InsertManyResult, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		topicsCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(VARIABLES_COLLECTION)
		insertResult, err := topicsCollection.InsertMany(context.TODO(), variables)
		if err != nil {
			slog.Error("fail to insert variables", "variables", variables, "error", err)
			return insertResult, err
		}
		return insertResult, nil
	}
	return nil, nil
}

func (d *DatabaseManager) insertVariable(variable *UserVariable) (*UserVariable, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		variablesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(VARIABLES_COLLECTION)
		result, err := variablesCollection.InsertOne(context.TODO(), variable)
		if err != nil {
			return nil, err
		}

		variableId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for node")
		}
		variable.Id = variableId
		return variable, nil
	}
	return nil, nil
}

func (d *DatabaseManager) updateVariable(variable *UserVariable) (*UserVariable, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		variablesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(VARIABLES_COLLECTION)
		filter := bson.D{{Key: "_id", Value: variable.Id}}
		update := bson.D{{Key: "$set", Value: variable}}
		_, err := variablesCollection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return nil, err
		}
		return variable, nil
	}
	return nil, nil
}

func (d *DatabaseManager) deleteVariable(id string) (*UserVariable, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {

		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		filter := bson.D{{Key: "_id", Value: objectId}}
		userVariables, err := d.findUserVariables(filter)
		if err != nil {
			return nil, err
		}

		if len(userVariables) == 0 {
			return nil, fmt.Errorf("variable with id: %s does not exist", id)
		}
		existingVariable := userVariables[0]

		variablesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(VARIABLES_COLLECTION)

		_, err = variablesCollection.DeleteOne(context.TODO(), filter)
		if err != nil {
			return nil, err
		}
		return &existingVariable, nil
	}
	return nil, nil
}
