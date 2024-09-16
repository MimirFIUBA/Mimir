package controllers

import (
	"encoding/json"
	"fmt"
	"mimir/internal/mimir"
	"mimir/triggers"
	"net/http"

	"github.com/gorilla/mux"
)

type TriggerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	SensorID    string `json:"sensorId"`
	TriggerType string `json:"type"`
	// Condition   triggers.Condition `json:"condition"`
	// Actions     []triggers.Action  `json:"actions"`
}

func GetTriggers(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Get Triggers")
	triggersBySensorId := mimir.Data.GetTriggersBySensorId()
	fmt.Printf("Triggers: %v \n", triggersBySensorId)
	var triggerResponseList []TriggerResponse

	for sensorId, triggerList := range triggersBySensorId {
		for _, trigger := range triggerList {
			var triggerResponse TriggerResponse
			switch trigger := trigger.(type) {
			case *triggers.Trigger:
				triggerResponse = TriggerResponse{"0", trigger.Name, sensorId, "trigger"}
			case *triggers.TimeTrigger:
				triggerResponse = TriggerResponse{"0", trigger.Name, sensorId, "time trigger"}
			default:
				panic("Bad Trigger!")
			}
			triggerResponseList = append(triggerResponseList, triggerResponse)
		}
	}

	fmt.Printf("triggers response: %v \n", triggerResponseList)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(triggerResponseList)
}

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sensor := mimir.Data.GetSensor(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func CreateTrigger(w http.ResponseWriter, r *http.Request) {

	var sensor *mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensor = mimir.Data.AddSensor(sensor)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sensor)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	var sensor *mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensor = mimir.Data.UpdateSensor(sensor)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mimir.Data.DeleteSensor(id)
	w.WriteHeader(http.StatusOK)
}
