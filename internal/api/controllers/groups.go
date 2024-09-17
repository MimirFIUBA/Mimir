package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	mimir "mimir/internal/mimir/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	groups := db.GroupsData.GetGroups()

	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   groups,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetGroupById(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	group, err := db.GroupsData.GetGroupById(id)
	if err != nil {
		logger.Error("Error searching for group", "group_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusNotFound, responses.GroupErrorCodes.NotFound)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   group,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var newGroup *mimir.Group
	err := json.NewDecoder(r.Body).Decode(&newGroup)
	if err != nil {
		logger.Error("Error decoding new group", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
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
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	if !db.GroupsData.IdExists(id) {
		logger.Error("Error updating group", "group_id", id, "error", "group doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.GroupErrorCodes.NotFound)
		return
	}

	var group *mimir.Group
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		logger.Error("Error decoding new group", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}
	group.ID = id

	group, err = db.GroupsData.UpdateGroup(group)
	if err != nil {
		logger.Error("Error updating group", "group_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.GroupErrorCodes.UpdateFailed)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected group was updated",
		Items:   group,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	if !db.GroupsData.IdExists(id) {
		logger.Error("Error deleting group", "group_id", id, "error", "group doesnt exist")
		responses.SendErrorResponse(w, http.StatusNotFound, responses.GroupErrorCodes.NotFound)
	}

	err := db.GroupsData.DeleteGroup(id)
	if err != nil {
		logger.Error("Error deleting group", "group_id", id, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.GroupErrorCodes.DeleteFailed)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The group was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}
