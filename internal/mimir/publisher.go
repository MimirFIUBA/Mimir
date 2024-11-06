package mimir

import (
	"context"
	"log/slog"
	"mimir/internal/models"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	client     mqtt.Client
	msgChannel chan models.MqttOutgoingMessage
}

func NewPublisher(client mqtt.Client, msgChannel chan models.MqttOutgoingMessage) *Publisher {
	return &Publisher{client, msgChannel}
}

func (p *Publisher) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case outgoingMessage := <-p.msgChannel:
			wg.Add(1)
			go func() {
				defer wg.Done()
				token := p.client.Publish(outgoingMessage.Topic, 0, false, outgoingMessage.Message)
				token.Wait()
				slog.Info("publish message to topic", "topic", outgoingMessage.Topic, "message", outgoingMessage.Message)
			}()
		case <-ctx.Done():
			slog.Error("context done, publisher", "error", ctx.Err())
			return
		}

	}
}
