package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddTriggersRoutes(router *mux.Router) {
	triggerRouter := router.PathPrefix("/triggers").Subrouter()

	triggerRouter.HandleFunc("", controllers.GetTriggers).Methods("GET")
	triggerRouter.HandleFunc("", controllers.CreateTrigger).Methods("POST")
	triggerRouter.HandleFunc("/{id}", controllers.GetTrigger).Methods("GET")
	triggerRouter.HandleFunc("/{id}", controllers.UpdateTrigger).Methods("PUT")
	triggerRouter.HandleFunc("/{id}", controllers.DeleteTrigger).Methods("DELETE")
}
