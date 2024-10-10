package mongodb

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Connect() (*mongo.Client, error) {

	dbUsername := os.Getenv("MONGODB_USERNAME")
	if dbUsername == "" {
		return nil, errors.New("MONGODB_USERNAME must be set")
	}

	dbPassword := os.Getenv("MONGODB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb+srv://%s:%s@mimir.razfo.mongodb.net/?retryWrites=true&w=majority&appName=Mimir", dbUsername, dbPassword)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to MongoDB!")
	return client, nil
}
