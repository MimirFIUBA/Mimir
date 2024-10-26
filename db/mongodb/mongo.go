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

	dbProtocol := os.Getenv("MONGODB_PROTOCOL")
	if dbProtocol == "" {
		return nil, errors.New("MONGODB_PROTOCOL must be set")
	}

	dbUsername := os.Getenv("MONGODB_USERNAME")
	if dbUsername == "" {
		return nil, errors.New("MONGODB_USERNAME must be set")
	}

	dbPassword := os.Getenv("MONGODB_PASSWORD")
	if dbPassword == "" {
		return nil, errors.New("MONGODB_PASSWORD must be set")
	}

	dbHostname := os.Getenv("MONGODB_HOSTNAME")
	if dbHostname == "" {
		return nil, errors.New("MONGODB_HOSTNAME must be set")
	}

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("%s://%s:%s@%s/?retryWrites=true&w=majority&appName=Mimir", dbProtocol, dbUsername, dbPassword, dbHostname)
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
	return client, nil
}
