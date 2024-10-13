package controllers

import (
	"encoding/binary"
	"encoding/json"
	"io"
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
			configuration, ok := configurationInterface.(map[string]interface{})
			if ok {
				pathInterface, exists := configuration["path"]
				if exists {
					path, ok := pathInterface.(string)
					if ok {
						jsonProcessor.AddValueConfiguration(processors.NewJSONValueConfiguration("", path))
					}
				}
			}
		}
		messageProcessor = jsonProcessor
	case "bytes":
		bytesProcessor := &processors.BytesProcessor{
			Name:            requestBody.Name,
			Topic:           requestBody.Topic,
			Type:            requestBody.ProcessorType,
			ReadingsChannel: MimirProcessor.ReadingChannel}
		//TODO implement
		for _, configurationInterface := range requestBody.Configurations {
			configuration, ok := configurationInterface.(map[string]interface{})
			if ok {
				endiannessInterface, exists := configuration["endianness"]
				if exists {
					endianness, ok := endiannessInterface.(string)
					if ok {
						bytesProcessor.AddBytesConfiguration(*processors.NewBytesConfiguration(endianness, binary.BigEndian, 1))
					}
				}
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(requestBody)
}

func UpdateProcessor(w http.ResponseWriter, r *http.Request) {
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
	w.WriteHeader(http.StatusOK)
}
