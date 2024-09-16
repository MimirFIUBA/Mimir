package routes

import (
	"fmt"
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	fmt.Println("router")
	router := mux.NewRouter()

	AddSensorRoutes(router)
	AddNodesRoutes(router)
	AddGroupRoutes(router)
	AddProcessorsRoutes(router)
	AddWebSocketRoutes(router)

	go controllers.HandleWebSocketMessages() //TODO (#26)

	return router
}
