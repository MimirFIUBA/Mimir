package controllers

import (
	"encoding/json"
	"io"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/internal/mimir/processors"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetProcessors(w http.ResponseWriter, r *http.Request) {
	processors := mimir.MessageProcessors.GetProcessors()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processors)
}

func GetProcessor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	processor, exists := mimir.MessageProcessors.GetProcessor(strings.ReplaceAll(id, ".", "/"))
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(processor)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateProcessor(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	decoder := json.NewDecoder(r.Body)
	var requestBody responses.Processor
	for {
		err := decoder.Decode(&requestBody)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	_, exists := mimir.MessageProcessors.GetProcessor(requestBody.Topic)
	if exists {
		logger.Error("Error creating new processor", "body", r.Body, "error", "processor for topic "+requestBody.Topic+" already exists")
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.AlreadyExists)
		return
	}

	var messageProcessor processors.MessageProcessor
	switch requestBody.ProcessorType {
	case "json":
		jsonProcessor := &processors.JSONProcessor{
			Name:            requestBody.Name,
			Topic:           requestBody.Topic,
			Type:            requestBody.ProcessorType,
			ReadingsChannel: MimirProcessor.ReadingChannel}
		//TODO see if we can do this better
		for _, configurationInterface := range requestBody.Configurations {
			jsonConfigurationMap, ok := configurationInterface.(map[string]interface{})
			if ok {
				jsonConfiguration, err := processors.JsonMapToJsonConfiguration(jsonConfigurationMap)
				if err != nil {
					logger.Error("Error creating new processor", "body", r.Body, "error", err)
					responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
					return
				}
				jsonProcessor.AddValueConfiguration(jsonConfiguration)
			}
		}
		messageProcessor = jsonProcessor
	case "bytes":
		bytesProcessor := &processors.BytesProcessor{
			Name:            requestBody.Name,
			Topic:           requestBody.Topic,
			Type:            requestBody.ProcessorType,
			ReadingsChannel: MimirProcessor.ReadingChannel}
		for _, configurationInterface := range requestBody.Configurations {
			configurationMap, ok := configurationInterface.(map[string]interface{})
			if ok {
				byteConfiguration, err := processors.JsonMapToByteConfiguration(configurationMap)
				if err != nil {
					logger.Error("Error creating new processor", "body", r.Body, "error", err)
					responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
					return
				}
				bytesProcessor.AddBytesConfiguration(*byteConfiguration)
			}
		}
		messageProcessor = bytesProcessor
	}

	_, err := db.SensorsData.GetSensorByTopic(requestBody.Topic)
	if err != nil {
		sensor := models.NewSensor(requestBody.Topic)
		sensor.Topic = requestBody.Topic
		MimirProcessor.RegisterSensor(sensor)
	}

	mimir.MessageProcessors.RegisterProcessor(requestBody.Topic, messageProcessor)

	db.Database.SaveProcessor(messageProcessor)

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The new processor was created",
		Items:   messageProcessor,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateProcessor(w http.ResponseWriter, r *http.Request) {
	//TODO: missing implementation
	var sensor *models.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func DeleteProcessor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	topic := strings.ReplaceAll(id, ".", "/")
	processor, exists := mimir.MessageProcessors.GetProcessor(topic)
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	db.Database.DeleteProcessor(processor)
	mimir.MessageProcessors.RemoveProcessor(topic)

	w.WriteHeader(http.StatusOK)
}
