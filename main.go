package main

import (
	"fmt"
	"mimir/internal/api/logging"
	"mimir/internal/api/routes"
	"net/http"
)

func main() {
	// TODO(27) - Replace hardcoded config with real config
	loggerConfig := logging.LoggerConfiguration{
		Level:    "debug",
		Filename: "api.log",
	}

	apiLogger, err := logging.CreateAPILogger(&loggerConfig)
	if err != nil {
		panic(err)
	}

	router := routes.CreateRouter(apiLogger)

	apiLogger.Info("Starting server at port 8080")
	// TODO(27) - Delete hardcoded port
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		panic(fmt.Sprintf("Error starting server: %s\n", err))
	}
}
