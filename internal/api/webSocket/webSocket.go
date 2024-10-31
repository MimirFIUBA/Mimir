package websocket

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/api/responses"
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	BroadcastChan chan string
	Clients       map[*websocket.Conn]bool
	Upgrader      websocket.Upgrader
}

func NewHandler() *WebSocketHandler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &WebSocketHandler{nil, make(map[*websocket.Conn]bool), upgrader}
}

func (h *WebSocketHandler) Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return h.Upgrader.Upgrade(w, r, nil)
}

func (h *WebSocketHandler) NewConnection(conn *websocket.Conn) {
	h.Clients[conn] = true
}

func (h *WebSocketHandler) CloseConnection(client *websocket.Conn) {
	client.Close()
	delete(h.Clients, client)
}

func (h *WebSocketHandler) BroadcastMessage(msg responses.WSMessage) {
	h.BroadcastChan <- fmt.Sprintf("%#v", msg)
}

// HandleWebSocketMessages listens to the broadcastChan and sends the message received from it to all clients.
func (h *WebSocketHandler) HandleMessages(ctx context.Context) {
	for {
		select {
		case msg := <-h.BroadcastChan:
			for client := range h.Clients {
				err := client.WriteJSON(msg)
				if err != nil {
					fmt.Println(err)
					h.CloseConnection(client)
				}
			}
		case <-ctx.Done():
			slog.Info("web socket context done, closing web socket handler", "error", ctx.Err())
			return
		}

	}
}
