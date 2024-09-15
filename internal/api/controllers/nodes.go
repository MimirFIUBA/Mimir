package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"net/http"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	nodes := db.NodesData.GetNodes()

	// TODO(#19) - Improve error handling
	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected nodes information was returned",
		Items:   nodes,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func GetNodeById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	node, err := db.NodesData.GetNodeById(id)
	if err != nil {
		logger.Error("Error searching for node", "node_id", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected node information was returned",
		Items:   node,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func CreateNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	var newNode *models.Node
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		logger.Error("Error decoding new node", "body", r.Body, "error", err.Error())
		return
	}

	_ = db.NodesData.CreateNode(newNode)
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new node was created",
		Items:   newNode,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func UpdateNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var node *models.Node
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		logger.Error("Error decoding new node", "body", r.Body, "error", err.Error())
		return
	}
	node.ID = id

	node, err = db.NodesData.UpdateNode(node)
	if err != nil {
		logger.Error("Error updating nodes", "node_id", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected node was updated",
		Items:   node,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.NodesData.DeleteNode(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		logger.Error("Error deleting group", "group_id", id, "error", err.Error())
		return
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The node was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}
