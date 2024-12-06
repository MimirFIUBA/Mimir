package controllers

import (
	"encoding/json"
	"mimir/internal/api/middlewares"
	"mimir/internal/api/responses"
	"mimir/internal/db"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserVariables(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())
	var items []db.UserVariable

	db.UserVariables.Range(func(key any, value any) bool {
		// name, ok := key.(string)
		// if ok {
		stringValue, ok := value.(*db.UserVariable)
		if ok {
			items = append(items, *stringValue)
		}
		// if ok {
		// 	items[name] = stringValue
		// }
		// }
		return true
	})

	err := responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
		Code:    0,
		Message: "All user variables information was returned",
		Items:   items,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func GetUserVariableByName(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	name := mux.Vars(r)["id"]
	userVariable, exists := db.GetUserVariable(name)

	var err error
	if !exists {
		err = responses.SendJSONResponse(w, http.StatusNotFound, responses.ItemsResponse{
			Code:    0,
			Message: "All user variables information was returned",
			Items:   nil,
		})
	} else {
		err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
			Code:    0,
			Message: "All user variables information was returned",
			Items:   userVariable,
		})
	}

	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
	}
}

func CreateUserVariable(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	var newVariable *db.UserVariable
	err := json.NewDecoder(r.Body).Decode(&newVariable)
	if err != nil {
		logger.Error("Error decoding new variable", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}

	userVariable, err := db.AddUserVariable(newVariable.Name, newVariable, true)
	if err != nil {
		logger.Error("Error inserting variable to database", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
	err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
		Code:    0,
		Message: "The new variable was created",
		Items:   userVariable,
	})
	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
		return
	}
}

func UpdateUserVariable(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	var newVariable *db.UserVariable
	err := json.NewDecoder(r.Body).Decode(&newVariable)
	if err != nil {
		logger.Error("Error decoding new variable", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}
	newVariable.Id, err = primitive.ObjectIDFromHex(id)
	if err != nil {
		logger.Error("Error decoding new variable", "body", r.Body, "error", err.Error())
		responses.SendErrorResponse(w, http.StatusBadRequest, responses.GroupErrorCodes.InvalidSchema)
		return
	}

	existingVariable, exists := db.GetUserVariableById(id)
	if exists {
		userVariable, err := db.UpdateUserVariable(existingVariable.Name, newVariable, true)
		if err != nil {
			logger.Error("Error updating variable", "error", err.Error())
			responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
			return
		}
		err = responses.SendJSONResponse(w, http.StatusCreated, responses.ItemsResponse{
			Code:    0,
			Message: "Variable updated successfully",
			Items:   userVariable,
		})
		if err != nil {
			logger.Error("Error sending response", "error", err.Error())
			responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
			return
		}
	} else {
		err = responses.SendJSONResponse(w, http.StatusNotFound, responses.ItemsResponse{
			Code:    0,
			Message: "Variable could not be found",
			Items:   nil,
		})
		if err != nil {
			logger.Error("Error sending response", "error", err.Error())
			responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
			return
		}

	}

}

func DeleteUserVariable(w http.ResponseWriter, r *http.Request) {
	logger := middlewares.ContextWithLogger(r.Context())

	id := mux.Vars(r)["id"]
	userVariable, err := db.DeleteUserVariable(id)

	if userVariable == nil {
		err = responses.SendJSONResponse(w, http.StatusNotFound, responses.ItemsResponse{
			Code:    0,
			Message: "User variable not found",
			Items:   nil,
		})
	} else {
		if err == nil {
			err = responses.SendJSONResponse(w, http.StatusOK, responses.ItemsResponse{
				Code:    0,
				Message: "User variable deleted succesfuly",
				Items:   userVariable,
			})
		}
	}

	if err != nil {
		logger.Error("Error sending response", "error", err.Error())
		responses.SendErrorResponse(w, http.StatusInternalServerError, responses.InternalErrorCodes.ResponseError)
	}
}
