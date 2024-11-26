package controllers

import (
	"encoding/json"
	"io"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"mimir/internal/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	sensors := db.SensorsData.GetSensors()
	nodes := db.NodesData.GetNodes()
	groups := db.GroupsData.GetGroups()

	items := make([]responses.SensorResponse, 0)

	groupsMap := make(map[string]*responses.GroupResponse)

	for _, group := range groups {
		groupsMap[group.GetId()] = responses.NewGroupResponse(*group)
	}

	nodesMap := make(map[string]*responses.NodeResponse)

	for _, node := range nodes {
		nodesMap[node.GetId()] = responses.NewNodeResponse(*node)
	}

	for _, sensor := range sensors {
		sensorResponse := responses.NewSensorResponse(*sensor)
		nodeForSensor, exists := nodesMap[sensor.NodeID]
		if exists {
			groupForNode, exists := groupsMap[nodeForSensor.GroupID]
			if exists {
				nodeForSensor.Group = *groupForNode
			}
			sensorResponse.Node = *nodeForSensor
		}
		items = append(items, *sensorResponse)
	}

	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   items,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetSensorById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
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

	newSensor, _ := CreateSensorInternal(w, r)

	err := responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new sensor was created",
		Items:   newSensor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
	}
}

func CreateSensorInternal(w http.ResponseWriter, r *http.Request) (*models.Sensor, string) {
	logger := middlewares.ContextWithLogger(r.Context())

	bytedata, err := io.ReadAll(r.Body)
	bodyString := string(bytedata)

	if err != nil {
		logger.Error("Error decoding new sensorHandler", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return nil, ""
	}

	var newSensor *models.Sensor
	err = json.Unmarshal([]byte(bodyString), &newSensor)

	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return nil, ""
	}

	_ = db.SensorsData.CreateSensor(newSensor)

	return newSensor, bodyString
}

func UpdateSensor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	if !db.SensorsData.IdExists(id) {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", "sensor doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
		return
	}

	var sensor *models.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		logger.Error("Error decoding new sensor", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return
	}

	sensor, err = db.SensorsData.UpdateSensor(sensor, id)
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
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	if !db.SensorsData.IdExists(id) {
		logger.Error("Error searching for sensors", "sensor_id", id, "error", "sensor doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
		return
	}

	err := db.SensorsData.DeleteSensor(id)
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
