package routes

import "github.com/gorilla/mux"

func CreateRouter() *mux.Router {
	router := mux.NewRouter()

	AddSensorRoutes(router)

	return router
}
