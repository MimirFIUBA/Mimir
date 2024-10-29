package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"mimir/internal/handlers"
	"mimir/internal/mimir"
	"mimir/internal/models"
	"mimir/internal/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetHandlers(w http.ResponseWriter, r *http.Request) {
	handlers := mimir.Mimir.MsgProcessor.GetHandlers()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handlers)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	handler, exists := mimir.Mimir.MsgProcessor.GetHandler(strings.ReplaceAll(id, ".", "/"))
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var requestBody responses.Handler
	err := utils.DecodeJsonToMap(r.Body, &requestBody)
	if err != nil {
		logger.Error("Error updating processor", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
	}

	_, exists := mimir.Mimir.MsgProcessor.GetHandler(requestBody.Topic)
	if exists {
		logger.Error("Error creating new processor", "body", r.Body, "error", "processor for topic "+requestBody.Topic+" already exists")
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.AlreadyExists)
		return
	}

	var messageHandler handlers.MessageHandler
	switch requestBody.HandlerType {
	case "json":
		jsonHandler, err := createJSONHandler(requestBody)
		if err != nil {
			logger.Error("Error creating new processor", "body", r.Body, "error", err)
			responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
		}
		messageHandler = jsonHandler
	case "bytes":
		bytesHandler, err := createBytesHandler(requestBody)
		if err != nil {
			logger.Error("Error creating new processor", "body", r.Body, "error", err)
			responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
		}
		messageHandler = bytesHandler
	}

	sensor, err := db.SensorsData.GetSensorByTopic(requestBody.Topic)
	if err != nil {
		sensor = models.NewSensor(requestBody.Topic)
		sensor.Topic = requestBody.Topic
	}
	MimirEngine.RegisterSensor(sensor)

	mimir.Mimir.MsgProcessor.RegisterHandler(requestBody.Topic, messageHandler)

	db.Database.SaveHandler(messageHandler)

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The new processor was created",
		Items:   messageHandler,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]
	topic := strings.ReplaceAll(id, ".", "/")

	var requestBody map[string]interface{}
	err := utils.DecodeJsonToMap(r.Body, &requestBody)
	if err != nil {
		logger.Error("Error updating processor", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
	}

	existingHandler, exists := mimir.Mimir.MsgProcessor.GetHandler(topic)
	if !exists {
		logger.Error("Error updating processor", "body", r.Body, "error", "processor for topic "+topic+" does not exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.ProcessorErrorCodes.NotFound)
		return
	}

	err = existingHandler.UpdateFields(requestBody)
	if err != nil {
		logger.Error("Error updating processor", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
		return
	}
	db.Database.SaveHandler(existingHandler)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingHandler)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	topic := strings.ReplaceAll(id, ".", "/")
	processor, exists := mimir.Mimir.MsgProcessor.GetHandler(topic)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db.Database.DeleteHandler(processor)
	mimir.Mimir.MsgProcessor.RemoveHandler(topic)

	w.WriteHeader(http.StatusOK)
}

// TODO: pasar esto a un handler factory con el reading channel para sacarlo de mimir engine
func createJSONHandler(requestBody responses.Handler) (*handlers.JSONHandler, error) {
	jsonHandler := &handlers.JSONHandler{
		Name:            requestBody.Name,
		Topic:           requestBody.Topic,
		Type:            requestBody.HandlerType,
		ReadingsChannel: MimirEngine.ReadingChannel}
	for _, configurationInterface := range requestBody.Configurations {
		jsonConfigurationMap, ok := configurationInterface.(map[string]interface{})
		if ok {
			jsonConfiguration, err := handlers.JsonMapToJsonConfiguration(jsonConfigurationMap)
			if err != nil {
				return nil, err
			}
			jsonHandler.AddValueConfiguration(jsonConfiguration)
		}
	}
	return jsonHandler, nil
}

// TODO: pasar esto a un handler factory con el reading channel para sacarlo de mimir engine
func createBytesHandler(requestBody responses.Handler) (*handlers.BytesHandler, error) {
	bytesHandler := &handlers.BytesHandler{
		Name:            requestBody.Name,
		Topic:           requestBody.Topic,
		Type:            requestBody.HandlerType,
		ReadingsChannel: MimirEngine.ReadingChannel}
	for _, configurationInterface := range requestBody.Configurations {
		configurationMap, ok := configurationInterface.(map[string]interface{})
		if ok {
			byteConfiguration, err := handlers.JsonMapToByteConfiguration(configurationMap)
			if err != nil {
				return nil, err
			}
			bytesHandler.AddBytesConfiguration(*byteConfiguration)
		}
	}
	return bytesHandler, nil
}
