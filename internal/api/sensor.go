package api

import (
	"encoding/json"
	mimir "mimir/internal/mimir"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type sensorsResponse struct {
	Sensors []mimir.Sensor `json:"sensors"`
}

func getSensors(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func getSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	key := vars["key"]

	id, err := strconv.Atoi(key)

	if err != nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	sensor := mimir.GetSensor(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func createSensor(w http.ResponseWriter, r *http.Request) {

	var sensor mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensor = mimir.CreateSensor(sensor)

	json.NewEncoder(w).Encode(sensor)
}

func updateSensor(w http.ResponseWriter, r *http.Request) {
	var sensor mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func deleteSensor(w http.ResponseWriter, _ *http.Request) {

	w.WriteHeader(http.StatusOK)
}
