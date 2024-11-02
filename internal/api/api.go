package api

import (
	"context"
	"log"
	"sync"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/mimir"
	"net/http"
)

func Start(ctx context.Context, wg *sync.WaitGroup, mimirEngine *mimir.MimirEngine) {
	controllers.SetMimirEngine(mimirEngine)
	controllers.Handler = controllers.NewWebSocketHandler(mimirEngine.WsChannel)
	router := routes.CreateRouter()
	go controllers.Handler.HandleMessages(ctx, wg)

	go func() {
		log.Fatal(http.ListenAndServe(":8080", router))
	}()
}
