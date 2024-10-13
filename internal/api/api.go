package api

import (
	"log"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/mimir"
	"net/http"
)

func Start(mimirProcessor *mimir.MimirProcessor) {
	router := routes.CreateRouter()
	controllers.SetWebSocketBroadcastChan(mimirProcessor.WsChannel)
	controllers.SetMimirProcessor(mimirProcessor)
	go controllers.HandleWebSocketMessages()
	log.Fatal(http.ListenAndServe(":8080", router))
}
