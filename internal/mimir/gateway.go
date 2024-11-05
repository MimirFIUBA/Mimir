package mimir

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"mimir/internal/handlers"
	"mimir/internal/models"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"golang.org/x/sync/errgroup"
)

type Gateway struct {
	id           string
	client       mqtt.Client
	done         chan struct{}
	readingsChan chan models.SensorReading
	timeout      time.Duration
	quiesce      uint
	retries      int
	qos          byte
	topics       *sync.Map
	msgs         MessageChannel
}

type GatewayOptions struct {
	ID      string
	Broker  string
	Timeout time.Duration
	Quiesce uint
	Retries int
	QoS     byte
}

func NewGateway(readingsChannel chan models.SensorReading, msgs MessageChannel, opts *GatewayOptions) (*Gateway, error) {
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
		readingsChan: readingsChannel,
		retries:      opts.Retries,
		done:         make(chan struct{}, 1),
		timeout:      opts.Timeout,
		quiesce:      opts.Quiesce,
		qos:          opts.QoS,
		topics:       new(sync.Map),
		msgs:         msgs,
	}, nil
}

func (g *Gateway) Start(topics <-chan []string, ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case newTopics := <-topics:
			wg.Add(1)
			go func() {
				defer wg.Done()
				topicsToSubscribe := make([]string, 0)
				for _, topic := range newTopics {
					isSubscribed, exists := g.topics.Load(topic)

					if !exists {
						topicsToSubscribe = append(topicsToSubscribe, topic)
					} else {
						isSubscribedBool, ok := isSubscribed.(bool)
						if !ok {
							slog.Error("wrong value topic subscribed", "topic", topic)
							continue
						}
						if !isSubscribedBool {
							topicsToSubscribe = append(topicsToSubscribe, topic)
						}
					}
				}

				g.trySubscribeToTopics(ctx, topicsToSubscribe)
			}()
		case <-ctx.Done():
			slog.Error("context done, gateway", "error", ctx.Err())
			return
		}
	}
}

func (g *Gateway) trySubscribeToTopics(ctx context.Context, topics []string) error {
	eg, _ := errgroup.WithContext(ctx)
	for _, topic := range topics {
		topic := topic
		eg.Go(func() error {
			token := g.client.Subscribe(topic, g.qos, onMessageReceived(g.msgs))
			var lastError error
			for i := 0; i < g.retries; i++ {
				timeout := time.Duration(math.Pow(2, float64(i))) * g.timeout
				if token.WaitTimeout(timeout) {
					if err := token.Error(); err != nil {
						lastError = err
						continue
					}
					g.topics.Store(topic, true)
					break
				}
			}
			if lastError != nil {
				return fmt.Errorf("couldn't subscribe to topic %q: %w", topic, lastError)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func onMessageReceived(ch MessageChannel) mqtt.MessageHandler {
	return func(c mqtt.Client, m mqtt.Message) {
		fmt.Println("Received message from: ", m.Topic())
		msg := handlers.Message{Topic: m.Topic(), Payload: m.Payload()}
		ch <- msg
		m.Ack() //TODO: check if the ack is necessary
	}

}

func tryConnectToBroker(ctx context.Context, client mqtt.Client) error {
	token := client.Connect()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-token.Done():
		return token.Error()
	}
}

func (g *Gateway) GetClient() mqtt.Client {
	return g.client
}

func (g *Gateway) CloseConnection() {
	topicsToUnsuscribe := make([]string, 0)
	g.topics.Range(func(key, value any) bool {
		isSubscribed, ok := value.(bool)
		if ok && isSubscribed {
			topic, ok := key.(string)
			if ok {
				topicsToUnsuscribe = append(topicsToUnsuscribe, topic)
			}
		}
		return true
	})
	var lastError error
	for i := 0; i < g.retries; i++ {
		timeout := time.Duration(math.Pow(2, float64(i))) * g.timeout
		token := g.client.Unsubscribe(topicsToUnsuscribe...)
		if token.WaitTimeout(timeout) {
			if err := token.Error(); err != nil {
				lastError = err
				continue
			}
			break
		}
	}

	if lastError != nil {
		slog.Error("could not unsuscribe to topics", "error", lastError)
	}

	g.client.Disconnect(g.quiesce)
}
