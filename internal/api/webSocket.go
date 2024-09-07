package api

import (
	"fmt"
	"net/http"
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	clients[conn] = true

	for {
		var msg WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			delete(clients, conn)
			return
		}

		fmt.Sprintf("%#v", msg)

		broadcast <- fmt.Sprintf("%#v", msg)
	}
}

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	ws, _ := upgrader.Upgrade(w, r, nil)
// 	defer ws.Close()

// 	for {
// 		// Receive message
// 		mt, message, _ := ws.ReadMessage()
// 		log.Printf("Message received: %s", message)

// 		// Reverse message
// 		n := len(message)
// 		for i := 0; i < n/2; i++ {
// 			message[i], message[n-1-i] = message[n-1-i], message[i]
// 		}

// 		// Response message
// 		_ = ws.WriteMessage(mt, message)
// 		_ = ws.WriteMessage(mt, message)
// 		_ = ws.WriteMessage(mt, message)
// 		log.Printf("Message sent: %s", message)
// 	}
// }

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
