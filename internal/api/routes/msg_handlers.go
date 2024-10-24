package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddHandlersRoutes(router *mux.Router) {
	processorsRouter := router.PathPrefix("/handlers").Subrouter()
	processorsRouter.HandleFunc("", controllers.GetHandlers).Methods("GET")
	processorsRouter.HandleFunc("", controllers.CreateHandler).Methods("POST")
	processorsRouter.HandleFunc("/{id}", controllers.GetHandler).Methods("GET")
	processorsRouter.HandleFunc("/{id}", controllers.UpdateHandler).Methods("PUT")
	processorsRouter.HandleFunc("/{id}", controllers.DeleteHandler).Methods("DELETE")
}
