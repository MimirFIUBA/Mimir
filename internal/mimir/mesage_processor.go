package mimir

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/handlers"
	"sync"
)

type MessageProcessor struct {
	handlers map[string]handlers.MessageHandler
	messages MessageChannel
}

func NewMessageProcessor(msgChannel MessageChannel) *MessageProcessor {
	return &MessageProcessor{handlers: make(map[string]handlers.MessageHandler), messages: msgChannel}
}

type MessageChannel chan handlers.Message

func (p *MessageProcessor) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case message := <-p.messages:
			fmt.Printf("New message from: %s - payload: %s\n", message.Topic, message.Payload)
			topic := message.Topic
			wg.Add(1)
			go func() {
				defer wg.Done()
				p.handlers[topic].HandleMessage(message)
			}()
		case <-ctx.Done():
			slog.Error("context done, message processor", "error", ctx.Err())
			return
		}
	}
}

func (p *MessageProcessor) RegisterHandler(topic string, handler handlers.MessageHandler) {
	p.handlers[topic] = handler
}

func (p *MessageProcessor) RemoveHandler(topic string) {
	delete(p.handlers, topic)
}

func (p *MessageProcessor) GetHandler(topic string) (handlers.MessageHandler, bool) {
	handler, exists := p.handlers[topic]
	return handler, exists
}
func (p *MessageProcessor) GetHandlers() map[string]handlers.MessageHandler {
	return p.handlers
}
