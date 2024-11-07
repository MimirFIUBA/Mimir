package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddUserVariablesRoutes(router *mux.Router) {
	userVariablesRouter := router.PathPrefix("/user-variables").Subrouter()

	userVariablesRouter.HandleFunc("", controllers.GetUserVariables).Methods("GET")
	userVariablesRouter.HandleFunc("", controllers.CreateUserVariable).Methods("POST")
	userVariablesRouter.HandleFunc("/{id}", controllers.GetUserVariableByName).Methods("GET")
	userVariablesRouter.HandleFunc("/{id}", controllers.UpdateUserVariable).Methods("PUT")
	userVariablesRouter.HandleFunc("/{id}", controllers.DeleteUserVariable).Methods("DELETE")
}
