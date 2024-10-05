package controllers

import (
	"encoding/json"
	"io"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/internal/utils"
	"net/http"

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

	processor, exists := mimir.MessageProcessors.GetProcessor(id)
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(processor)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateProcessor(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonMap map[string]interface{}
	for {
		err := decoder.Decode(&jsonMap)
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	topic, ok := jsonMap["topic"].(string)
	if !ok {
		http.Error(w, badRequestErrorMessage, http.StatusBadRequest)
		return
	}

	processorType, ok := jsonMap["type"].(string)
	if !ok {
		http.Error(w, badRequestErrorMessage, http.StatusBadRequest)
		return
	}

	processor, err := utils.JsonToProcessor(processorType, jsonMap)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mimir.MessageProcessors.RegisterProcessor(topic, processor)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(processor)
}

func UpdateProcessor(w http.ResponseWriter, r *http.Request) {
	var sensor *models.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sensor = mimir.Data.UpdateSensor(sensor)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func DeleteProcessor(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// id := vars["id"]

	// mimir.Data.DeleteSensor(id)
	w.WriteHeader(http.StatusOK)
}