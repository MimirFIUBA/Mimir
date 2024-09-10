package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/models"
	"mimir/internal/api/utils"
	"net/http"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	sensors := db.SensorsData.GetSensors()

	// TODO(#19) - Improve error handling
	err := utils.RespondWithJSONItems(w, http.StatusOK, sensors)
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
	err = utils.RespondWithJSONItems(w, http.StatusOK, sensor)
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
	err = utils.RespondWithJSONItems(w, http.StatusCreated, newSensor)
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
	err = utils.RespondWithJSONItems(w, http.StatusOK, sensor)
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.SensorsData.DeleteSensor(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		fmt.Printf("Error deleting sensor: %s", err)
	}

	// TODO - Change response
	err = utils.RespondWithJSONItems(w, http.StatusNoContent, nil)
}
