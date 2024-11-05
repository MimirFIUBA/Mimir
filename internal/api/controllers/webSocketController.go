package controllers

import (
	"context"
	"fmt"
	"mimir/internal/api/responses"
	websocket "mimir/internal/api/webSocket"
	"net/http"
)

var WSHandler = websocket.NewHandler()

func SetWebSocketBroadcastChan(broadcastChan chan string) {
	WSHandler.BroadcastChan = broadcastChan
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := WSHandler.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	WSHandler.NewConnection(conn)

	// TODO: solo para testing, broadcasteo los mensajes que recibo
	for {
		var msg responses.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			WSHandler.CloseConnection(conn)
			return
		}
		WSHandler.BroadcastMessage(msg)
	}
}

func HandleWebSocketMessages(ctx context.Context) {
	WSHandler.HandleMessages(ctx)
}
