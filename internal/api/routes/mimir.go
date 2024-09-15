package routes

import (
	"github.com/gorilla/mux"
	"log/slog"
	"mimir/internal/api/middlewares"
)

func CreateRouter(apiLogger *slog.Logger, requestLogger *slog.Logger) *mux.Router {
	router := mux.NewRouter()

	// Adds Middlewares
	router.Use(middlewares.CreateAPILoggerMiddleware(apiLogger))
	router.Use(middlewares.CreateRequestLoggerMiddleware(requestLogger))

	// Adds Routes
	AddSensorRoutes(router)
	AddNodesRoutes(router)
	AddGroupRoutes(router)

	return router
}
