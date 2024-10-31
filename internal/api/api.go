package api

import (
	"context"
	"log"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/mimir"
	"net/http"
)

func Start(ctx context.Context, mimirEngine *mimir.MimirEngine) {
	router := routes.CreateRouter()
	controllers.SetWebSocketBroadcastChan(mimirEngine.WsChannel)
	controllers.SetMimirEngine(mimirEngine)
	go controllers.HandleWebSocketMessages(ctx)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()
}
