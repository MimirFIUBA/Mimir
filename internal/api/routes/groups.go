package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddGroupRoutes(router *mux.Router) *mux.Router {
	groupsRouter := router.PathPrefix("/groups").Subrouter()

	groupsRouter.HandleFunc("/", controllers.GetGroups).Methods("GET")
	groupsRouter.HandleFunc("/", controllers.CreateGroup).Methods("POST")
	groupsRouter.HandleFunc("/{id}", controllers.GetGroupById).Methods("GET")
	groupsRouter.HandleFunc("/{id}", controllers.UpdateGroup).Methods("PUT")
	groupsRouter.HandleFunc("/{id}", controllers.DeleteGroup).Methods("DELETE")

	return groupsRouter
}
