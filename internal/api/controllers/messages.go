package controllers

import (
	"encoding/json"
	"fmt"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/models"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	messages, err := db.Database.FindAllMessages()
	if err != nil {
		logger.Error("Error finding messages", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.GroupErrorCodes.UpdateFailed)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All messages returned",
		Items:   messages,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	handler, exists := mimir.Mimir.MsgProcessor.GetHandler(strings.ReplaceAll(id, ".", "/"))
	if exists {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(handler)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var newMessage *models.Message
	err := json.NewDecoder(r.Body).Decode(&newMessage)
	if err != nil {
		logger.Error("Error decoding message body", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return
	}

	newMessage, err = db.Database.InsertMessage(newMessage)
	if err != nil {
		logger.Error("Error inserting message", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.SensorErrorCodes.InvalidSchema)
		return
	}
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new sensor was created",
		Items:   newMessage,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}

}

func UpdateMessage(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	var messageUpdate *models.Message
	err := json.NewDecoder(r.Body).Decode(&messageUpdate)

	if err != nil {
		logger.Error("Error decoding new variable", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}

	updatedMessage, err := db.Database.UpdateMessage(id, messageUpdate)
	if err != nil {
		logger.Error("Error updating variable", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}

	err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "The message was updated",
		Items:   updatedMessage,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}

}

func DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Println(id)

	w.WriteHeader(http.StatusOK)
}
