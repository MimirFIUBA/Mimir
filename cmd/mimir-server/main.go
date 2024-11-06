package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log/slog"
	"mimir/internal/api"
	"mimir/internal/config"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"net"
	"os"
	"os/signal"
	"strings"
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

func reloadTriggers(mimirEngine *mimir.MimirEngine) {
	fmt.Println("Reload triggers")
	config.BuildTriggers(mimirEngine)
}

func reloadHandlers(mimirEngine *mimir.MimirEngine) {
	fmt.Println("Reload handlers")
	config.BuildHandlers(mimirEngine)
}

func processCommand(command string, mimirEngine *mimir.MimirEngine) {
	switch command {
	case "reloadTriggers":
		reloadTriggers(mimirEngine)
	case "reloadHandlers":
		reloadHandlers(mimirEngine)
	default:
		fmt.Println("bad command")
	}
}

func cli(mimirEngine *mimir.MimirEngine) {
	// Escuchar en el puerto 8080 para conexiones TCP
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		fmt.Println("Error al crear el servidor TCP:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Servidor TCP escuchando en el puerto 8082")

	for {
		// Aceptar conexiones entrantes
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error al aceptar conexión:", err)
			continue
		}

		// Manejar cada conexión en una goroutine
		go func(c net.Conn) {
			defer c.Close()
			reader := bufio.NewReader(c)
			for {
				input, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					fmt.Println("Error al leer comando:", err)
					return
				}
				command := strings.TrimSpace(input)
				processCommand(command, mimirEngine)
			}
		}(conn)
	}
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

	go cli(mimirEngine)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	close(sigChan)

	slog.Info("closing application")

	mimirEngine.Close()

	cancel()
	wg.Wait()
	api.Stop()
	fmt.Println("Waiting for all processes to finsih")
	slog.Info("close successful")

	fmt.Println("Mimir is out of duty, bye!")
}
