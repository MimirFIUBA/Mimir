package api

import (
	"encoding/json"
	mimir "mimir/internal/mimir"
	"net/http"

	"github.com/gorilla/mux"
)

type sensorsResponse struct {
	Sensors []mimir.Sensor `json:"sensors"`
}

func getSensors(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.Data.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func getSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	sensor := mimir.Data.GetSensor(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensor)
}

func createSensor(w http.ResponseWriter, r *http.Request) {

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

func updateSensor(w http.ResponseWriter, r *http.Request) {
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

func deleteSensor(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mimir.Data.DeleteSensor(id)
	w.WriteHeader(http.StatusOK)
}
