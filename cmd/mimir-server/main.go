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
	"syscall"
)

func main() {

	handlerOpts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, handlerOpts))
	slog.SetDefault(logger)

	slog.Info("Starting")
	fmt.Println("MiMiR starting")

	config.LoadIni()
	slog.Info("ini config file loaded")

	// mimirProcessor := mimir.NewMimirProcessor()
	mimirProcessor := mimir.StartMimir()
	mimirProcessor.StartGateway()
	slog.Info("gateway started")

	mongoClient, err := db.Database.ConnectToMongo()
	if err != nil {
		slog.Error("error connecting to mongo", "error", err)
	} else {
		slog.Info("connection to mongo succesfully established")
		defer func() {
			if err = mongoClient.Disconnect(context.TODO()); err != nil {
				panic(err)
			}
		}()
	}

	influxClient, err := db.Database.ConnectToInfluxDB()
	if err != nil {
		slog.Error("error connecting to influx db", "error", err)
	} else {
		slog.Info("connection to influxdb succesfully established")
		defer influxClient.Close()
	}

	config.BuildInitialConfiguration(mimirProcessor)
	slog.Info("succesfully built environment based on configuration")
	db.Run()

	go mimirProcessor.Run()
	go api.Start(mimirProcessor)

	fmt.Println("Everything up and running")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	slog.Info("closing application")
	mimir.CloseConnection()
	slog.Info("close successful")

	fmt.Println("Mimir is out of duty, bye!")
}
