package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/config"
	"mimir/internal/db"
	"mimir/internal/utils"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetTriggers(w http.ResponseWriter, _ *http.Request) {

	triggers := db.Database.GetTriggers()

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

	trigger := config.BuildTriggerObserver(*newTrigger, MimirProcessor)
	db.RegisterTrigger(trigger, newTrigger.Topics)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTrigger)
}

func UpdateTrigger(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	vars := mux.Vars(r)
	id := vars["id"]

	var requestBody db.Trigger
	err := utils.DecodeJsonToMap(r.Body, &requestBody)
	if err != nil {
		logger.Error("Error updating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
		return
	}

	updatedTrigger, err := db.Database.UpdateTrigger(id, &requestBody)
	if err != nil {
		logger.Error("Error updating trigger", "body", r.Body, "error", err)
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.ProcessorErrorCodes.InvalidSchema)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTrigger)
}

func DeleteTrigger(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("not implemented")
}
