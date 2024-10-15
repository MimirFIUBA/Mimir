package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()

	AddSensorRoutes(router)
	AddNodesRoutes(router)
	AddGroupRoutes(router)
	AddProcessorsRoutes(router)
	AddTriggersRoutes(router)
	AddWebSocketRoutes(router)

	go controllers.HandleWebSocketMessages() //TODO (#26)

	return router
}
