package mimir

import (
	"context"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/factories"
	"mimir/internal/models"
	"mimir/triggers"
	"sync"

	"github.com/gookit/ini/v2"
)

type MimirEngine struct {
	ReadingChannel chan models.SensorReading
	TopicChannel   chan []string
	WsChannel      chan models.WSOutgoingMessage
	ActionFactory  *factories.ActionFactory
	TriggerFactory *factories.TriggerFactory
	MsgProcessor   *MessageProcessor
	gateway        *Gateway
	publisher      *Publisher
	firstCancel    context.CancelFunc
	secondCancel   context.CancelFunc
	firstWg        *sync.WaitGroup
	secondWg       *sync.WaitGroup
	Scheduler      *triggers.Scheduler
}

func NewMimirEngine() *MimirEngine {
	topicChannel := make(chan []string)
	readingsChannel := make(chan models.SensorReading, 50)
	outgoingMessagesChannel := make(chan models.MqttOutgoingMessage)
	webSocketMessageChannel := make(chan models.WSOutgoingMessage)
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

	scheduler, err := triggers.NewScheduler()
	if err != nil {
		slog.Error("Fail to create scheduler", "error", err)
	}

	engine := MimirEngine{
		ReadingChannel: readingsChannel,
		TopicChannel:   topicChannel,
		WsChannel:      webSocketMessageChannel,
		ActionFactory:  factories.NewActionFactory(outgoingMessagesChannel, webSocketMessageChannel),
		TriggerFactory: factories.NewTriggerFactory(),
		MsgProcessor:   NewMessageProcessor(msgChannel),
		gateway:        gateway,
		publisher:      NewPublisher(gateway.GetClient(), outgoingMessagesChannel),
		Scheduler:      scheduler,
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
	e.Scheduler.Start()
}

func (e *MimirEngine) processReadings(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case reading := <-e.ReadingChannel:
			slog.Info("new reading", "reading", reading)
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
	e.Scheduler.Shutdown()
	e.firstCancel()
	e.firstWg.Wait()
	e.secondCancel()
	e.secondWg.Wait()

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
	slog.Info("Registering sensor", "topic", sensor.Topic)
	sensor.IsActive = true
	e.TopicChannel <- []string{sensor.Topic}
}

func (e *MimirEngine) BuildTrigger(trigger factories.TriggerOptions) (triggers.Trigger, error) {
	return e.TriggerFactory.BuildTrigger(trigger)
}
