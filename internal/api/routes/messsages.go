package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddMessagesRoutes(router *mux.Router) {
	messagesRouter := router.PathPrefix("/messages").Subrouter()
	messagesRouter.HandleFunc("", controllers.GetMessages).Methods("GET")
	messagesRouter.HandleFunc("", controllers.CreateMessage).Methods("POST")
	messagesRouter.HandleFunc("/{id}", controllers.GetMessage).Methods("GET")
	messagesRouter.HandleFunc("/{id}", controllers.UpdateMessage).Methods("PUT")
	messagesRouter.HandleFunc("/{id}", controllers.DeleteMessage).Methods("DELETE")
}
