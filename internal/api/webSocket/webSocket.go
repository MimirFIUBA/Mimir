package websocket

import (
	"fmt"
	"mimir/internal/api/responses"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	BroadcastChan chan string
	Clients       map[*websocket.Conn]bool
	Upgrader      websocket.Upgrader
}

func NewHandler() *Handler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &Handler{nil, make(map[*websocket.Conn]bool), upgrader}
}

func (h *Handler) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return h.Upgrader.Upgrade(w, r, nil)
}

func (h *Handler) NewConnection(conn *websocket.Conn) {
	h.Clients[conn] = true
}

func (h *Handler) CloseConnection(client *websocket.Conn) {
	client.Close()
	delete(h.Clients, client)
}

func (h *Handler) BroadcastMessage(msg responses.WSMessage) {
	h.BroadcastChan <- fmt.Sprintf("%#v", msg)
}

// HandleWebSocketMessages listens to the broadcastChan and sends the message received from it to all clients.
func (h *Handler) HandleMessages() {
	for {
		msg := <-h.BroadcastChan

		for client := range h.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				h.CloseConnection(client)
			}
		}
	}
}
