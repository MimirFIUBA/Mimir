package api

import (
	"encoding/json"
	mimir "mimir/internal/mimir"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func getNodes(w http.ResponseWriter, r *http.Request) {
	nodes := mimir.Data.GetNodes()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}

func getNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	node := mimir.Data.GetNode(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(node)
}

func createNode(w http.ResponseWriter, r *http.Request) {

	var node *mimir.Node
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node = mimir.Data.AddNode(node)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(node)
}

func updateNode(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}

func deleteNode(w http.ResponseWriter, r *http.Request) {
	var sensors = sensorsResponse{
		Sensors: mimir.GetSensors(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(sensors)
}
