package mimir

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/models"
	"mimir/triggers"
	"sync"

	"github.com/gookit/ini/v2"
)

type MimirEngine struct {
	ReadingChannel chan models.SensorReading
	TopicChannel   chan []string
	WsChannel      chan string
	ActionFactory  *models.ActionFactory
	TriggerFactory *models.TriggerFactory
	MsgProcessor   *MessageProcessor
	gateway        *Gateway
	publisher      *Publisher
	firstCancel    context.CancelFunc
	secondCancel   context.CancelFunc
	firstWg        *sync.WaitGroup
	secondWg       *sync.WaitGroup
}

func NewMimirEngine() *MimirEngine {
	topicChannel := make(chan []string)
	readingsChannel := make(chan models.SensorReading, 50)
	outgoingMessagesChannel := make(chan string)
	webSocketMessageChannel := make(chan string)
	msgChannel := make(MessageChannel)

	opts := GatewayOptions{
		ID:      "1",
		Broker:  ini.String(consts.MQTT_BROKER_CONFIG_NAME),
		Timeout: consts.MQTT_SUBSCRIPTION_TIMEOUT,
		Quiesce: consts.MQTT_QUIESCE,
		Retries: consts.MQTT_MAX_RETRIES,
		QoS:     consts.MQTT_QOS,
	}

	gateway, err := NewGateway(readingsChannel, msgChannel, &opts)
	if err != nil {
		slog.Error("Error creating new gateway", "error", err)
	}

	engine := MimirEngine{
		ReadingChannel: readingsChannel,
		TopicChannel:   topicChannel,
		WsChannel:      webSocketMessageChannel,
		ActionFactory:  models.NewActionFactory(outgoingMessagesChannel, webSocketMessageChannel),
		TriggerFactory: models.NewTriggerFactory(),
		MsgProcessor:   NewMessageProcessor(msgChannel),
		gateway:        gateway,
		publisher:      NewPublisher(gateway.GetClient(), outgoingMessagesChannel),
	}
	Mimir = &engine
	return &engine
}

func (e *MimirEngine) Run(ctx context.Context) {
	publisherCtx, publisherCancel := context.WithCancel(ctx)
	generalCtx, cancel := context.WithCancel(publisherCtx)
	e.firstCancel = cancel
	e.secondCancel = publisherCancel
	e.firstWg = &sync.WaitGroup{}
	e.secondWg = &sync.WaitGroup{}

	go e.publisher.Run(publisherCtx, e.secondWg)
	go e.MsgProcessor.Run(generalCtx, e.firstWg)
	go e.processReadings(generalCtx, e.firstWg)
	go e.gateway.Start(e.TopicChannel, generalCtx, e.firstWg)
}

func (e *MimirEngine) processReadings(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case reading := <-e.ReadingChannel:
			wg.Add(1)
			go func() {
				defer wg.Done()
				db.StoreReading(reading)
			}()
		case <-ctx.Done():
			slog.Info("context done, processReadings", "error", ctx.Err())
			return
		}
	}
}

func (e *MimirEngine) Close() {
	fmt.Println("First cancel")
	e.firstCancel()
	e.firstWg.Wait()

	fmt.Println("Second cancel")
	e.secondCancel()
	e.secondWg.Wait()
	fmt.Println("finish wait")

	e.gateway.CloseConnection()
	setTriggersInactive()
	setTopicsInactive()
}

func setTopicsInactive() {
	db.SensorsData.SetSensorsToInactive()
	db.Database.DeactivateTopics(db.SensorsData.GetSensors())
}

func setTriggersInactive() {
	db.Database.DeactivateTriggers(context.TODO())
}

// TODO: ver de sacar esto de aca
func (e *MimirEngine) RegisterSensor(sensor *models.Sensor) {
	fmt.Println("Register sensor ", sensor.Topic)
	sensor.IsActive = true
	e.TopicChannel <- []string{sensor.Topic}
}

func (e *MimirEngine) BuildTrigger(trigger models.TriggerOptions) (triggers.Trigger, error) {
	return e.TriggerFactory.BuildTrigger(trigger)
}
