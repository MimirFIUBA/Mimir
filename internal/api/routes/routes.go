package routes

import (
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()

	AddSensorRoutes(router)
	AddNodesRoutes(router)
	AddGroupRoutes(router)
	AddHandlersRoutes(router)
	AddTriggersRoutes(router)
	AddWebSocketRoutes(router)

	return router
}
