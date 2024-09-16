package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddSensorRoutes(router *mux.Router) {
	sensorRouter := router.PathPrefix("/sensors").Subrouter()

	sensorRouter.HandleFunc("", controllers.GetSensors).Methods("GET")
	sensorRouter.HandleFunc("", controllers.CreateSensor).Methods("POST")
	sensorRouter.HandleFunc("/{id}", controllers.GetSensorById).Methods("GET")
	sensorRouter.HandleFunc("/{id}", controllers.UpdateSensor).Methods("PUT")
	sensorRouter.HandleFunc("/{id}", controllers.DeleteSensor).Methods("DELETE")
}
