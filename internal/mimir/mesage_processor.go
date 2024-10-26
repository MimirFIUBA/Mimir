package mimir

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/handlers"
)

type MessageProcessor struct {
	handlers map[string]handlers.MessageHandler
	messages MessageChannel
}

func NewMessageProcessor(msgChannel MessageChannel) *MessageProcessor {
	return &MessageProcessor{handlers: make(map[string]handlers.MessageHandler), messages: msgChannel}
}

type MessageChannel chan handlers.Message

func (p *MessageProcessor) Run(ctx context.Context) {
	for {
		select {
		case message := <-p.messages:
			fmt.Println("New message from ", message.Topic)
			topic := message.Topic
			go p.handlers[topic].HandleMessage(message)
		case <-ctx.Done():
			slog.Error("context", "error", ctx.Err())
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
