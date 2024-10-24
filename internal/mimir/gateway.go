package mimir

import (
	"context"
	"fmt"
	"mimir/internal/models"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Gateway struct {
	id           string
	client       mqtt.Client
	done         chan struct{}
	readingsChan chan models.SensorReading
	topicsChan   chan string
	timeout      time.Duration
	quiesce      uint
	retries      int
	qos          byte
}

type GatewayOptions struct {
	ID      string
	Broker  string
	Timeout time.Duration
	Quiesce uint
	Retries int
	QoS     byte
}

func NewGateway(msgs chan models.SensorReading, opts *GatewayOptions) (*Gateway, error) {
	mqttOptions := mqtt.NewClientOptions()
	mqttOptions.AddBroker(opts.Broker)
	mqttOptions.SetProtocolVersion(4)
	client := mqtt.NewClient(mqttOptions)
	deadlineCtx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()
	if err := tryConnectToBroker(deadlineCtx, client); err != nil {
		return nil, fmt.Errorf("couldn't connect to broker: client %s: %w", opts.ID, deadlineCtx.Err())
	}
	return &Gateway{
		id:           opts.ID,
		client:       client,
		readingsChan: msgs,
		retries:      opts.Retries,
		done:         make(chan struct{}, 1),
		timeout:      opts.Timeout,
		quiesce:      opts.Quiesce,
		qos:          opts.QoS,
	}, nil
}

// func (g *Gateway) Start(topics <-chan []string, ctx context.Context) {
// 	var subscribed []string
// 	for {
// 		select {
// 		case newTopics := <-topics:
// 			newSet := SetFromSlice(newTopics)
// 			new, deleted := newSet.GetNewAndDeletedTopics(subscribed)
// 			subscribed = newTopics
// 			// NOTE(juan): Try a number of times if it fails, the
// 			// correct way would be to do an exponential
// 			// backoff
// 			for i := 0; i < g.retries; i++ {
// 				if err := trySubscribeToTopics(ctx, new, g); err != nil {
// 					slog.Error("subscribe", "error", err)
// 					continue
// 				}
// 				break
// 			}
// 			for i := 0; i < g.retries; i++ {
// 				if err := tryUnsubscribeToTopics(deleted, g); err != nil {
// 					slog.Error("unsubscribe", "error", err)
// 					continue
// 				}
// 				break
// 			}
// 		case <-ctx.Done():
// 			slog.Error("context", "error", ctx.Err())
// 			return
// 		}
// 	}
// }

// func onMessageReceived2(ch MessageChannel) mqtt.MessageHandler {
// 	return func(c mqtt.Client, m mqtt.Message) {
// 		payload := m.Payload()
// 		topic := m.Topic()
// 		msg := Message{topic, payload}
// 		ch <- msg
// 		m.Ack()
// 		return
// 	}
// }

// func trySubscribeToTopics(ctx context.Context, topics []string, g *Gateway) error {
// 	eg, ctx := errgroup.WithContext(ctx)
// 	for _, topic := range topics {
// 		topic := topic
// 		eg.Go(func() error {
// 			token := g.client.Subscribe(topic, g.qos, onMessageReceived2(g.msgs))
// 			if token.WaitTimeout(g.timeout) {
// 				if err := token.Error(); err != nil {
// 					return fmt.Errorf("couldn't subscribe to topic %q: %w", topic, err)
// 				}
// 			}
// 			return nil
// 		})
// 	}
// 	if err := eg.Wait(); err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (g *Gateway) Done() <-chan struct{} {
// 	return g.done
// }

// func (g *Gateway) Shutdown() {
// 	g.client.Disconnect(g.quiesce)
// }

func tryConnectToBroker(ctx context.Context, client mqtt.Client) error {
	token := client.Connect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		return token.Error()
	}
}
