package main

import (
	"flag"
	"fmt"
	"mimir/internal/api/configuration"
	"mimir/internal/api/logging"
	"mimir/internal/api/routes"
	"net/http"
)

func main() {
	configFile := flag.String("config", "configuration.json", "Path to the configuration file")
	flag.Parse()

	config, err := configuration.GetConfiguration(*configFile)
	if err != nil {
		panic(err)
	}

	apiLogger, err := logging.CreateAPILogger(&config.Logging)
	if err != nil {
		panic(err)
	}

	router := routes.CreateRouter(apiLogger)

	apiLogger.Info(fmt.Sprintf("Starting server at port %d", config.Server.Port))
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), router)
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %s\n", err))
	}
}
