package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/config"
	"mimir/internal/db"
	"mimir/internal/utils"
	"mimir/triggers"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTriggers(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	triggers, err := db.Database.GetAllTriggers()
	if err != nil {
		logger.Error("error getting all triggers", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(triggers)
}

func GetTrigger(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	vars := mux.Vars(r)
	id := vars["id"]

	trigger, err := db.Database.GetTrigger(id)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			responses.SendErrorResponse(w, http.StatusNotFound, responses.TriggerErrorCodes.NotFound)
			return
		}
		logger.Error("Error getting trigger", "error", err)
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.UnexpectedError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(trigger)
}

func CreateTrigger(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var requestBody db.Trigger
	err := utils.DecodeJsonToMap(r.Body, &requestBody)
	if err != nil {
		logger.Error("Error creating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
	}

	newTrigger, err := db.Database.InsertTrigger(&requestBody)
	if err != nil {
		logger.Error("Error creating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.UnexpectedError)
		return
	}

	trigger, err := config.BuildTrigger(*newTrigger, MimirEngine)
	if err != nil {
		if err.Error() == "condition does not compile" {
			logger.Error("Error creating trigger", "body", r.Body, "error", err)
			responses.SendErrorResponse(w, http.StatusBadRequest, responses.TriggerErrorCodes.ConditionDoesNotCompile)
			return
		}
		logger.Error("Error creating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.TriggerErrorCodes.UpdateFailed)
		return
	}

	db.RegisterTrigger(trigger, newTrigger.Topics)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrigger)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]

	//TODO: con este approach necesitamos mandar el cuerpo entero del trigger
	var requestBody db.Trigger
	err := utils.DecodeJsonToMap(r.Body, &requestBody)
	if err != nil {
		logger.Error("Error updating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.TriggerErrorCodes.InvalidSchema)
		return
	}

	actions := make([]triggers.Action, 0)
	for _, action := range requestBody.Actions {
		triggerAction := config.ToTriggerAction(action)
		actions = append(actions, triggerAction)
	}
	updatedTrigger, err := db.Database.UpdateTrigger(id, &requestBody, actions)
	if err != nil {
		logger.Error("Error updating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.TriggerErrorCodes.InvalidSchema)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrigger)
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]

	err := db.Database.DeleteTrigger(id)
	if err != nil {
		logger.Error("Error deleting trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.TriggerErrorCodes.InvalidSchema)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.MessageResponse{
		Code:    200,
		Message: "The trigger was deleted",
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}
