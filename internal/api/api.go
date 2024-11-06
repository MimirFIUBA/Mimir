package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"mimir/internal/api/controllers"
	"mimir/internal/api/routes"
	"mimir/internal/consts"
	"mimir/internal/mimir"
	"net/http"

	"github.com/gookit/ini/v2"
)

type ServerAddKey string

type ApiManager struct {
}

func Start(ctx context.Context, mimirEngine *mimir.MimirEngine) {
	router := routes.CreateRouter()
	controllers.SetWebSocketBroadcastChan(mimirEngine.WsChannel)
	controllers.SetMimirEngine(mimirEngine)

	addr := ":" + ini.String(consts.API_PORT_CONFIG_NAME)
	server := &http.Server{
		Addr:    addr,
		Handler: router,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, ServerAddKey("serverAddr"), l.Addr().String())
			return ctx
		},
	}

	go controllers.HandleWebSocketMessages(ctx)
	go func() {
		err := server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
	}()
}

func Stop() {
	controllers.WSHandler.Stop()
	slog.Info("Websocket stopped")
}
