package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"net/http"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	sensors := db.SensorsData.GetSensors()

	// TODO(#19) - Improve error handling
	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensors,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func GetSensorById(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	sensor, err := db.SensorsData.GetSensorById(id)
	if err != nil {
		fmt.Printf("Error searchinf for sensor with id %s: %s", id, err)
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensor,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func CreateSensor(w http.ResponseWriter, r *http.Request) {
	var newSensor *models.Sensor
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newSensor)
	if err != nil {
		fmt.Printf("Error decoding new sensor: %s", err)
		return
	}

	_ = db.SensorsData.CreateSensor(newSensor)
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new sensor was created",
		Items:   newSensor,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var sensor *models.Sensor
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		fmt.Printf("Error decoding body: %s", err)
		return
	}
	sensor.ID = id

	// TODO(#19) - Improve error handling
	sensor, err = db.SensorsData.UpdateSensor(sensor)
	if err != nil {
		fmt.Printf("Error updating sensor: %s", err)
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected sensor was updated",
		Items:   sensor,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.SensorsData.DeleteSensor(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		fmt.Printf("Error deleting sensor: %s", err)
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The sensor was deleted",
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}
