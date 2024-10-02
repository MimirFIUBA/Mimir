package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"mimir/internal/api/validators"
	"net/http"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	sensors := db.SensorsData.GetSensors()

	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensors,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetSensorById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidId)
		return
	}

	sensor, err := db.SensorsData.GetSensorById(id)
	if err != nil {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func CreateSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var newSensor *models.Sensor
	err := json.NewDecoder(r.Body).Decode(&newSensor)
	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
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
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidId)
		return
	}

	if !db.SensorsData.IdExists(id) {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", "sensor doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
		return
	}

	var sensor *models.Sensor
	err = json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return
	}
	sensor.ID = id

	sensor, err = db.SensorsData.UpdateSensor(sensor)
	if err != nil {
		logger.Error("Error updating sensor", "sensor_ud", id, "error", err.Error())
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected sensor was updated",
		Items:   sensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func DeleteSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidId)
		return
	}

	if !db.SensorsData.IdExists(id) {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", "sensor doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
		return
	}

	err = db.SensorsData.DeleteSensor(id)
	if err != nil {
		logger.Error("Error deleting sensor", "sensor_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.SensorErrorCodes.DeleteFailed)
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The sensor was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}
