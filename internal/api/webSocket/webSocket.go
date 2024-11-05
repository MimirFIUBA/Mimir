package websocket

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/api/responses"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketHandler struct {
	BroadcastChan chan string
	Clients       map[*websocket.Conn]bool
	Upgrader      websocket.Upgrader
	wg            sync.WaitGroup
}

func NewHandler() *WebSocketHandler {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &WebSocketHandler{
		Clients:  make(map[*websocket.Conn]bool),
		Upgrader: upgrader,
	}
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
				h.wg.Add(1)
				go func(client *websocket.Conn) {
					defer h.wg.Done()
					err := client.WriteJSON(msg)
					if err != nil {
						slog.Error("Error on websocket client", "error", err, "client", client)
						h.CloseConnection(client)
					}
				}(client)
			}
		case <-ctx.Done():
			slog.Info("web socket context done, closing web socket handler", "error", ctx.Err())
			return
		}
	}
}

func (h *WebSocketHandler) closeClients() {
	for client, isConnected := range h.Clients {
		if isConnected {
			err := client.Close()
			if err != nil {
				slog.Error("Error closing web socket client connection", "error", err, "client", client)
			}
		}
	}
	slog.Info("web socket clients closed")
}

func (h *WebSocketHandler) Stop() {
	h.wg.Wait()
	h.closeClients()
}
