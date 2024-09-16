package api

// import (
// 	"fmt"
// 	"mimir/internal/api/responses"
// 	"net/http"
// )

// func handleConnections(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	clients[conn] = true

// 	for {
// 		var msg responses.WSMessage
// 		err := conn.ReadJSON(&msg)
// 		if err != nil {
// 			fmt.Println(err)
// 			delete(clients, conn)
// 			return
// 		}
// 		broadcast <- fmt.Sprintf("%#v", msg)
// 	}
// }

// func handleMessages() {
// 	for {
// 		msg := <-broadcast

// 		for client := range clients {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				fmt.Println(err)
// 				client.Close()
// 				delete(clients, client)
// 			}
// 		}
// 	}
// }
