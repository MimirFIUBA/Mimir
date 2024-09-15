package main

import (
	"fmt"
	"mimir/internal/api/logging"
	"mimir/internal/api/routes"
	"net/http"
)

func main() {
	apiLogger := logging.CreateAPILogger()
	requestLogger := logging.CreateRequestLogger()

	router := routes.CreateRouter(apiLogger, requestLogger)

	fmt.Printf("Starting server at port 8080\n")
	// TODO - Delete hardcoded port
	err := http.ListenAndServe(":8080", router)
	// TODO - Improve error handling
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}
