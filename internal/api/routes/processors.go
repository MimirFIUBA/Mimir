package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddProcessorsRoutes(router *mux.Router) {
	processorsRouter := router.PathPrefix("/processors").Subrouter()
	processorsRouter.HandleFunc("", controllers.GetProcessors).Methods("GET")
	processorsRouter.HandleFunc("", controllers.CreateProcessor).Methods("POST")
	processorsRouter.HandleFunc("/{id}", controllers.GetProcessor).Methods("GET")
	processorsRouter.HandleFunc("/{id}", controllers.UpdateProcessor).Methods("PUT")
	processorsRouter.HandleFunc("/{id}", controllers.DeleteProcessor).Methods("DELETE")
}
