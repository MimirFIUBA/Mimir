package api

import (
	"encoding/json"
	"fmt"
	"mimir/internal/mimir"
	"net/http"

	"github.com/gorilla/mux"
)

func getGroups(w http.ResponseWriter, _ *http.Request) {
	var groups = mimir.Data.GetGroups()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(groups)
}

func getGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Get group - Id: %s\n", id)

	group := mimir.Data.GetGroup(id)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
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
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Update group - Id: %s\n", id)

	var group *mimir.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	group.ID = id

	group = mimir.Data.UpdateGroup(group)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(group)
}

func deleteGroup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("Delete group - Id: %s\n", id)

	mimir.Data.DeleteGroup(id)
	w.WriteHeader(http.StatusOK)
}
