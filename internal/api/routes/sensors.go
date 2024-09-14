package routes

import (
	"github.com/gorilla/mux"
	"mimir/internal/api/controllers"
)

func AddSensorRoutes(router *mux.Router) {
	sensorRouter := router.PathPrefix("/sensors").Subrouter()

	sensorRouter.HandleFunc("", controllers.GetSensors).Methods("GET")
	sensorRouter.HandleFunc("", controllers.CreateSensor).Methods("POST")
	sensorRouter.HandleFunc("/{id}", controllers.GetSensorById).Methods("GET")
	sensorRouter.HandleFunc("/{id}", controllers.UpdateSensor).Methods("PUT")
	sensorRouter.HandleFunc("/{id}", controllers.DeleteSensor).Methods("DELETE")
}
