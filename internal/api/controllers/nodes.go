package controllers

import (
	"encoding/json"
	"fmt"
	"mimir/internal/api/db"
	"mimir/internal/api/models"
	"mimir/internal/api/utils"
	"net/http"

	"github.com/gorilla/mux"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	nodes := db.NodesData.GetNodes()

	// TODO(#19) - Improve error handling
	err := utils.RespondWithJSONItems(w, http.StatusOK, nodes)
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func GetNodeById(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	node, err := db.NodesData.GetNodeById(id)
	if err != nil {
		fmt.Printf("Error searching for node with id %s: %s", id, err)
		return
	}

	// TODO(#19) - Improve error handling
	err = utils.RespondWithJSONItems(w, http.StatusOK, node)
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func CreateNode(w http.ResponseWriter, r *http.Request) {
	var newNode *models.Node
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		fmt.Printf("Error decoding new node: %s", err)
		return
	}

	_ = db.NodesData.CreateNode(newNode)
	err = utils.RespondWithJSONItems(w, http.StatusCreated, newNode)
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func UpdateNode(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var node *models.Node
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		fmt.Printf("Error decoding body: %s", err)
		return
	}
	node.ID = id

	node, err = db.NodesData.UpdateNode(node)
	if err != nil {
		fmt.Printf("Error updating node with id %s: %s", id, err)
		return
	}

	// TODO(#19) - Improve error handling
	err = utils.RespondWithJSONItems(w, http.StatusOK, node)
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.NodesData.DeleteNode(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		fmt.Printf("Error deleting node: %s", err)
	}

	// TODO - Change response
	err = utils.RespondWithJSONItems(w, http.StatusNoContent, nil)
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}
