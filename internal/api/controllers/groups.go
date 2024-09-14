package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"mimir/internal/api/db"
	"mimir/internal/api/models"
	"mimir/internal/api/responses"
	"net/http"
)

func GetGroups(w http.ResponseWriter, r *http.Request) {
	groups := db.GroupsData.GetGroups()

	// TODO(#19) - Improve error handling
	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   groups,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func GetGroupById(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	// TODO(#19) - Improve error handling
	group, err := db.GroupsData.GetGroupById(id)
	if err != nil {
		fmt.Printf("Error searching for group with id %s: %s", id, err)
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All selected groups information was returned",
		Items:   group,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func CreateGroup(w http.ResponseWriter, r *http.Request) {
	var newGroup *models.Group
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&newGroup)
	if err != nil {
		fmt.Printf("Error decoding new group: %s", err)
		return
	}

	_ = db.GroupsData.CreateGroup(newGroup)
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new group was created",
		Items:   newGroup,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]

	var group *models.Group
	// TODO(#19) - Improve error handling
	err := json.NewDecoder(r.Body).Decode(&group)
	if err != nil {
		fmt.Printf("Error decoding body: %s", err)
		return
	}
	group.ID = id

	// TODO(#19) - Improve error handling
	group, err = db.GroupsData.UpdateGroup(group)
	if err != nil {
		fmt.Printf("Error updating group: %s", err)
		return
	}

	// TODO(#19) - Improve error handling
	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The selected group was updated",
		Items:   group,
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	// TODO(#20) - Validate Query Params
	id := mux.Vars(r)["id"]
	err := db.GroupsData.DeleteGroup(id)

	// TODO(#19) - Improve error handling
	if err != nil {
		fmt.Printf("Error deleting group: %s", err)
	}

	err = responses.SendJSONResponse(w, http.StatusNoContent, responses.MessageResponse{
		Code:    0,
		Message: "The group was deleted",
	})
	if err != nil {
		fmt.Printf("Error responding with %s", err)
		return
	}
}
