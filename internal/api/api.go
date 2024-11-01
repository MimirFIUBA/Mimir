package api

import (
	"context"
	"log"
	"sync"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/api/websocket"
	"mimir/internal/mimir"
	"net/http"
)

func Start(ctx context.Context, wg *sync.WaitGroup, mimirEngine *mimir.MimirEngine) {
	controllers.SetMimirEngine(mimirEngine)
	controllers.WebSocketHandler = websocket.NewHandler(mimirEngine.WsChannel)
	router := routes.CreateRouter()
	go controllers.WebSocketHandler.HandleMessages(ctx, wg)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()
}
