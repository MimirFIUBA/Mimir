package mimir

import (
	"context"
	"log/slog"
	"mimir/internal/consts"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Publisher struct {
	client     mqtt.Client
	msgChannel chan string
}

func NewPublisher(client mqtt.Client, msgChannel chan string) *Publisher {
	return &Publisher{client, msgChannel}
}

func (p *Publisher) Run(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case outgoingMessage := <-p.msgChannel:
			wg.Add(1)
			go func() {
				defer wg.Done()
				topic := consts.MQTT_ALERT_TOPIC
				token := p.client.Publish(topic, 0, false, outgoingMessage)
				token.Wait()
				slog.Info("publish message to topic", "topic", topic, "message", outgoingMessage)
			}()
		case <-ctx.Done():
			slog.Error("context done, publisher", "error", ctx.Err())
			return
		}

	}
}
