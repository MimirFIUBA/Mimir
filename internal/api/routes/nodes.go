package routes

import (
	"github.com/gorilla/mux"
	"mimir/internal/api/controllers"
)

func AddNodesRoutes(router *mux.Router) {
	nodesRouter := router.PathPrefix("/nodes").Subrouter()

	nodesRouter.HandleFunc("/", controllers.GetNodes).Methods("GET")
	nodesRouter.HandleFunc("/", controllers.CreateNode).Methods("POST")
	nodesRouter.HandleFunc("/{id}", controllers.GetNodeById).Methods("GET")
	nodesRouter.HandleFunc("/{id}", controllers.UpdateNode).Methods("PUT")
	nodesRouter.HandleFunc("/{id}", controllers.DeleteNode).Methods("DELETE")
}
