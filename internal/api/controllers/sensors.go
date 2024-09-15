package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"net/http"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	sensors := db.SensorsData.GetSensors()

	// TODO(#19) - Improve error handling
	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensors,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func GetSensorById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	sensor, err := db.SensorsData.GetSensorById(id)
	if err != nil {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func CreateSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	var newSensor *models.Sensor
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newSensor)
	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		return
	}

	_ = db.SensorsData.CreateSensor(newSensor)
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new sensor was created",
		Items:   newSensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var sensor *models.Sensor
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		return
	}
	sensor.ID = id

	// TODO(#19) - Improve error handling
	sensor, err = db.SensorsData.UpdateSensor(sensor)
	if err != nil {
		logger.Error("Error updating sensor", "sensor_ud", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected sensor was updated",
		Items:   sensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.SensorsData.DeleteSensor(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		logger.Error("Error deleting sensor", "sensor_id", id, "error", err.Error())
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The sensor was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}
