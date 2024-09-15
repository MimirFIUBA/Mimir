package main

import (
	"fmt"
	"mimir/internal/api/logging"
	"mimir/internal/api/routes"
	"net/http"
)

func main() {
	apiLogger := logging.CreateAPILogger()

	router := routes.CreateRouter(apiLogger)

	fmt.Printf("Starting server at port 8080\n")
	// TODO - Delete hardcoded port
	err := http.ListenAndServe(":8080", router)
	// TODO - Improve error handling
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		return
	}
}
