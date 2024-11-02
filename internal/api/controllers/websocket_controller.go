package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/api/responses"
	"net/http"
	"sync"

	gws "github.com/gorilla/websocket"
)

var (
	Handler *WebSocketHandler
)

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := Handler.Upgrade(w, r)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	Handler.NewConnection(conn)

	// TODO: solo para testing, broadcasteo los mensajes que recibo
	for {
		var msg responses.WSMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println(err)
			Handler.CloseConnection(conn)
			return
		}
		Handler.BroadcastMessage(msg)
	}
}

type WebSocketHandler struct {
	broadcastChan chan string
	clients       map[*gws.Conn]bool
	upgrader      gws.Upgrader
}

func NewWebSocketHandler(broadcastChan chan string) *WebSocketHandler {
	upgrader := gws.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &WebSocketHandler{broadcastChan, make(map[*gws.Conn]bool), upgrader}
}

func (h *WebSocketHandler) Upgrade(w http.ResponseWriter, r *http.Request) (*gws.Conn, error) {
	return h.upgrader.Upgrade(w, r, nil)
}

func (h *WebSocketHandler) NewConnection(conn *gws.Conn) {
	h.clients[conn] = true
}

func (h *WebSocketHandler) CloseConnection(client *gws.Conn) {
	client.Close()
	delete(h.clients, client)
}

func (h *WebSocketHandler) BroadcastMessage(msg responses.WSMessage) {
	h.broadcastChan <- fmt.Sprintf("%#v", msg)
}

// HandleWebSocketMessages listens to the broadcastChan and sends the message received from it to all clients.
func (h *WebSocketHandler) HandleMessages(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case msg := <-h.broadcastChan:
			wg.Add(1)
			go func() {
				defer wg.Done()
				for client := range h.clients {
					err := client.WriteJSON(msg)
					if err != nil {
						fmt.Println(err)
						h.CloseConnection(client)
					}
				}
			}()
		case <-ctx.Done():
			slog.Info("web socket context done, closing web socket handler", "error", ctx.Err())
			return
		}

	}
}
