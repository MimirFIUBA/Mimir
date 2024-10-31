package main

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/api"
	"mimir/internal/config"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"os"
	"os/signal"
	"sync"
	"syscall"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func initializeDatabase() (*mongo.Client, influxdb2.Client) {
	mongoClient, err := db.Database.ConnectToMongo()
	if err != nil {
		slog.Error("error connecting to mongo", "error", err)
	} else {
		slog.Info("connection to mongo succesfully established")
	}

	influxClient, err := db.Database.ConnectToInfluxDB()
	if err != nil {
		slog.Error("error connecting to influx db", "error", err)
	} else {
		slog.Info("connection to influxdb succesfully established")
	}
	return mongoClient, influxClient
}

func initializeLogger() {
	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)
}

func loadStoredData(e *mimir.MimirEngine) {
	db.LoadTopology()
	config.BuildInitialConfiguration(e)
}

func gracefulShutdown(cancel context.CancelFunc, wg *sync.WaitGroup, e *mimir.MimirEngine) {
	slog.Info("closing application")

	cancel()
	wg.Wait()
	e.Close()
	fmt.Println("Waiting for all processes to finsih")
	slog.Info("close successful")

	fmt.Println("Mimir is out of duty, bye!")
}

func main() {

	initializeLogger()

	slog.Info("Starting")
	fmt.Println("Mimir starting")

	config.LoadIni()
	slog.Info("ini config file loaded")

	mongoClient, influxClient := initializeDatabase()
	if mongoClient != nil {
		defer func() {
			if err := mongoClient.Disconnect(context.TODO()); err != nil {
				slog.Error("error disconnecting from mongo", "error", err)
			}
		}()
	}

	if influxClient != nil {
		defer influxClient.Close()
	}

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	mimirEngine := mimir.NewMimirEngine()
	mimirEngine.Run(ctx)

	loadStoredData(mimirEngine)

	db.Run(ctx, &wg)
	fmt.Println("DB RUNNING")
	api.Start(ctx, mimirEngine)
	fmt.Println("API STARTED")

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	close(sigChan)

	slog.Info("closing application")

	mimirEngine.Close()

	cancel()
	wg.Wait()
	fmt.Println("Waiting for all processes to finsih")
	slog.Info("close successful")

	fmt.Println("Mimir is out of duty, bye!")
}
