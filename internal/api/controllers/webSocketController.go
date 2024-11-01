package controllers

import (
	"fmt"
	"mimir/internal/api/responses"
	"mimir/internal/api/websocket"
	"net/http"
)

var (
	WebSocketHandler *websocket.WebSocketHandler
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := WebSocketHandler.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	WebSocketHandler.NewConnection(conn)

	// TODO: solo para testing, broadcasteo los mensajes que recibo
	for {
		var msg responses.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			WebSocketHandler.CloseConnection(conn)
			return
		}
		WebSocketHandler.BroadcastMessage(msg)
	}
}
