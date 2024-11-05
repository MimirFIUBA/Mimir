package api

import (
	"context"
	"log"
	"log/slog"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/mimir"
	"net/http"
)

type ApiManager struct {
}

func Start(ctx context.Context, mimirEngine *mimir.MimirEngine) {
	router := routes.CreateRouter()
	controllers.SetWebSocketBroadcastChan(mimirEngine.WsChannel)
	controllers.SetMimirEngine(mimirEngine)
	go controllers.HandleWebSocketMessages(ctx)
	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()
}

func Stop() {
	controllers.WSHandler.Stop()
	slog.Info("Websocket stopped")
}
