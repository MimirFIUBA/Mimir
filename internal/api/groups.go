package api

import (
	"encoding/json"
	mimir "mimir/internal/mimir"
	"net/http"
)

func getGroups(w http.ResponseWriter, _ *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func createGroup(w http.ResponseWriter, r *http.Request) {

	var sensor mimir.Sensor
	err := json.NewDecoder(r.Body).Decode(&sensor)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sensor = mimir.CreateSensor(sensor)

	json.NewEncoder(w).Encode(sensor)
}

func updateGroup(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}
