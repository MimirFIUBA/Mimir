package api

import (
	"encoding/json"
	mimir "mimir/internal/mimir"
	"net/http"
)

func getGroups(w http.ResponseWriter, _ *http.Request) {
	var groups = mimir.Data.GetGroups()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func getGroup(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nil)
}

func createGroup(w http.ResponseWriter, r *http.Request) {

	var group *mimir.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group = mimir.Data.AddGroup(group)

	json.NewEncoder(w).Encode(group)
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
