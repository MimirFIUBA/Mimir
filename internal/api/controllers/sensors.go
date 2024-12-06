package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSensors(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	sensorsMap := db.SensorsData.GetSensorsMap()

	dbSensors, err := db.Database.FindAllTopics()
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}

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

	for _, sensor := range dbSensors {
		sensorResponse := responses.NewSensorResponse(sensor)
		memorySensor, exists := sensorsMap[sensor.Topic]
		if exists {
			sensorResponse.LastSensedReading = memorySensor.LastSensedReading
		}
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

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
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

	id := mux.Vars(r)["id"]
	sensor, err := db.SensorsData.GetSensorById(id)

	if err != nil {
		objectId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			logger.Error("Error searching for sensors", "sensor_id", id, "error", err.Error())
			responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
			return
		}
		sensors, err := db.Database.FindTopics(bson.D{{Key: "_id", Value: objectId}})
		if err != nil {
			logger.Error("Error searching for sensors", "sensor_id", id, "error", err.Error())
			responses.SendErrorResponse(w, http.StatusInternalServerError, responses.SensorErrorCodes.NotFound)
			return
		}
		if len(sensors) < 1 {
			logger.Error("Sensor not found", "sensor_id", id)
			responses.SendErrorResponse(w, http.StatusNotFound, responses.SensorErrorCodes.NotFound)
			return
		}
		sensor = &sensors[0]
	}

	sensorResponse := responses.NewSensorResponse(*sensor)

	handler, exists := mimir.Mimir.MsgProcessor.GetHandler(strings.ReplaceAll(sensor.Topic, ".", "/"))
	if exists {
		sensorResponse.Handler = &handler
	}

	triggers, err := db.Database.GetTriggers(bson.D{{Key: "topics", Value: sensor.Topic}})
	if err != nil {
		logger.Error("Error getting triggers", "error", err.Error())
	}

	sensorResponse.Triggers = triggers

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected sensors information was returned",
		Items:   sensorResponse,
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

	if newSensor.Topic == "" {
		node, err := db.NodesData.GetNodeById(newSensor.NodeID)
		if err == nil {
			group, err := db.GroupsData.GetGroupById(node.GroupID)
			if err == nil {
				newSensor.Topic = consts.MQTT_TOPIC_PREFIX + strings.ToLower(group.Name) + "/" + strings.ToLower(node.Name) + "/" + newSensor.DataName
			}
		}
	}

	createdSensor, err := db.SensorsData.CreateSensor(newSensor)
	if err != nil {
		logger.Error("Error creating sensor", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}

	trigger, err := MimirEngine.TriggerFactory.BuildNewReadingNotificationTrigger()
	if err != nil {
		logger.Error("Error creating reading notification trigger", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	} else {
		trigger.Activate()
		createdSensor.Register(trigger)
	}

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

	err = responses.SendJSONResponse(w, http.StatusOK, responses.MessageResponse{
		Code:    200,
		Message: "The sensor was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}
