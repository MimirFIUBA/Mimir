package controllers

import (
	"context"
	"fmt"
	"mimir/internal/api/responses"
	websocket "mimir/internal/api/webSocket"
	"net/http"
)

var handler = websocket.NewHandler()

func SetWebSocketBroadcastChan(broadcastChan chan string) {
	handler.BroadcastChan = broadcastChan
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := handler.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	handler.NewConnection(conn)

	// TODO: solo para testing, broadcasteo los mensajes que recibo
	for {
		var msg responses.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			handler.CloseConnection(conn)
			return
		}
		handler.BroadcastMessage(msg)
	}
}

func HandleWebSocketMessages(ctx context.Context) {
	handler.HandleMessages(ctx)
}
