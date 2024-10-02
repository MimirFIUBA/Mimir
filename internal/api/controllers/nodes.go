package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"mimir/internal/api/validators"
	"net/http"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	nodes := db.NodesData.GetNodes()

	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected nodes information was returned",
		Items:   nodes,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetNodeById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.NodeErrorCodes.InvalidId)
		return
	}

	node, err := db.NodesData.GetNodeById(id)
	if err != nil {
		logger.Error("Error searching for node", "node_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusNotFound, responses.NodeErrorCodes.NotFound)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected node information was returned",
		Items:   node,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func CreateNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var newNode *models.Node
	err := json.NewDecoder(r.Body).Decode(&newNode)
	if err != nil {
		logger.Error("Error decoding new node", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.NodeErrorCodes.InvalidSchema)
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
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.NodeErrorCodes.InvalidId)
		return
	}

	if !db.GroupsData.IdExists(id) {
		logger.Error("Error updating group", "node_id", id, "error", "group doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.NodeErrorCodes.NotFound)
		return
	}

	var node *models.Node
	err = json.NewDecoder(r.Body).Decode(&node)
	if err != nil {
		logger.Error("Error decoding new node", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.NodeErrorCodes.InvalidSchema)
		return
	}
	node.ID = id

	node, err = db.NodesData.UpdateNode(node)
	if err != nil {
		logger.Error("Error updating nodes", "node_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.NodeErrorCodes.UpdateFailed)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected node was updated",
		Items:   node,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func DeleteNode(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	err := validators.CheckIdNotEmpty(id)
	if err != nil {
		logger.Error("Error during id validation", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.NodeErrorCodes.InvalidId)
		return
	}

	if !db.GroupsData.IdExists(id) {
		logger.Error("Error deleting group", "node_id", id, "error", "group doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.NodeErrorCodes.NotFound)
		return
	}

	err = db.NodesData.DeleteNode(id)
	if err != nil {
		logger.Error("Error deleting group", "group_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.NodeErrorCodes.DeleteFailed)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The node was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}
