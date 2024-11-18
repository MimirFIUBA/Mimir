package controllers

import (
	"context"
	"log/slog"
	"mimir/internal/api/responses"
	websocket "mimir/internal/api/webSocket"
	"mimir/internal/models"
	"net/http"
)

var WSHandler = websocket.NewHandler()

func SetWebSocketBroadcastChan(broadcastChan chan models.WSOutgoingMessage) {
	WSHandler.BroadcastChan = broadcastChan
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := WSHandler.Upgrade(w, r)
	if err != nil {
		slog.Error("error on websocket upgrade", "error", err)
		return
	}
	defer conn.Close()

	WSHandler.NewConnection(conn)

	// TODO: solo para testing, broadcasteo los mensajes que recibo
	for {
		var msg responses.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			slog.Error("error reading json from websocket client", "error", err)
			WSHandler.CloseConnection(conn)
			return
		}
		WSHandler.BroadcastMessage(msg)
	}
}

func HandleWebSocketMessages(ctx context.Context) {
	WSHandler.HandleMessages(ctx)
}
