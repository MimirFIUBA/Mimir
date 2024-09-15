package routes

import (
	"mimir/internal/api/controllers"

	"github.com/gorilla/mux"
)

func AddWebSocketRoutes(router *mux.Router) {
	router.HandleFunc("/ws", controllers.HandleConnections)
}
