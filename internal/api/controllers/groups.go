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

func GetGroups(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	groups := db.GroupsData.GetGroups()

	// TODO(#19) - Improve error handling
	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   groups,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func GetGroupById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	group, err := db.GroupsData.GetGroupById(id)
	if err != nil {
		logger.Error("Error searching for group", "group_id", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   group,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	var newGroup *models.Group
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newGroup)
	if err != nil {
		logger.Error("Error decoding new group", "body", r.Body, "error", err.Error())
		return
	}

	_ = db.GroupsData.CreateGroup(newGroup)
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new group was created",
		Items:   newGroup,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var group *models.Group
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		logger.Error("Error decoding new group", "body", r.Body, "error", err.Error())
		return
	}
	group.ID = id

	// TODO(#19) - Improve error handling
	group, err = db.GroupsData.UpdateGroup(group)
	if err != nil {
		logger.Error("Error updating group", "group_id", id, "error", err.Error())
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected group was updated",
		Items:   group,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.GroupsData.DeleteGroup(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		logger.Error("Error deleting group", "group_id", id, "error", err.Error())
		return
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The group was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		return
	}
}
