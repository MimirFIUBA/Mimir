package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/processors"
	"mimir/internal/models"
	"mimir/internal/utils"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetHandlers(w http.ResponseWriter, r *http.Request) {
	handlers := mimir.MessageProcessors.GetHandlers()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(handlers)
}

func GetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	handler, exists := mimir.MessageProcessors.GetHandler(strings.ReplaceAll(id, ".", "/"))
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

	_, exists := mimir.MessageProcessors.GetHandler(requestBody.Topic)
	if exists {
		logger.Error("Error creating new processor", "body", r.Body, "error", "processor for topic "+requestBody.Topic+" already exists")
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.AlreadyExists)
		return
	}

	var messageHandler processors.MessageHandler
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

	_, err = db.SensorsData.GetSensorByTopic(requestBody.Topic)
	if err != nil {
		sensor := models.NewSensor(requestBody.Topic)
		sensor.Topic = requestBody.Topic
		MimirProcessor.RegisterSensor(sensor)
	}

	mimir.MessageProcessors.RegisterHandler(requestBody.Topic, messageHandler)

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

	existingHandler, exists := mimir.MessageProcessors.GetHandler(topic)
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
	processor, exists := mimir.MessageProcessors.GetHandler(topic)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db.Database.DeleteHandler(processor)
	mimir.MessageProcessors.RemoveHandler(topic)

	w.WriteHeader(http.StatusOK)
}

func createJSONHandler(requestBody responses.Handler) (*processors.JSONHandler, error) {
	jsonHandler := &processors.JSONHandler{
		Name:            requestBody.Name,
		Topic:           requestBody.Topic,
		Type:            requestBody.HandlerType,
		ReadingsChannel: MimirProcessor.ReadingChannel}
	for _, configurationInterface := range requestBody.Configurations {
		jsonConfigurationMap, ok := configurationInterface.(map[string]interface{})
		if ok {
			jsonConfiguration, err := processors.JsonMapToJsonConfiguration(jsonConfigurationMap)
			if err != nil {
				return nil, err
			}
			jsonHandler.AddValueConfiguration(jsonConfiguration)
		}
	}
	return jsonHandler, nil
}

func createBytesHandler(requestBody responses.Handler) (*processors.BytesHandler, error) {
	bytesHandler := &processors.BytesHandler{
		Name:            requestBody.Name,
		Topic:           requestBody.Topic,
		Type:            requestBody.HandlerType,
		ReadingsChannel: MimirProcessor.ReadingChannel}
	for _, configurationInterface := range requestBody.Configurations {
		configurationMap, ok := configurationInterface.(map[string]interface{})
		if ok {
			byteConfiguration, err := processors.JsonMapToByteConfiguration(configurationMap)
			if err != nil {
				return nil, err
			}
			bytesHandler.AddBytesConfiguration(*byteConfiguration)
		}
	}
	return bytesHandler, nil
}
