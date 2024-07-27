package api

import (
	"encoding/json"
	"fmt"
	mimir "mimir/internal/mimir"
	"net/http"

	"github.com/gorilla/mux"
)

func getNodes(w http.ResponseWriter, r *http.Request) {
	nodes := mimir.Data.GetNodes()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodes)
}

func getNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

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
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Update node - Id: %s\n", id)

	var node *mimir.Node
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	node.ID = id

	node = mimir.Data.UpdateNode(node)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(node)
}

func deleteNode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mimir.Data.DeleteNode(id)
	w.WriteHeader(http.StatusOK)
}
