package mimir

import (
	"context"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/models"
	"mimir/triggers"

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
		readingsChannel,
		topicChannel,
		webSocketMessageChannel,
		models.NewActionFactory(outgoingMessagesChannel, webSocketMessageChannel),
		models.NewTriggerFactory(),
		NewMessageProcessor(msgChannel),
		gateway,
		NewPublisher(gateway.GetClient(), outgoingMessagesChannel),
	}
	Mimir = &engine
	return &engine
}

func (e *MimirEngine) Run(ctx context.Context) {
	go e.gateway.Start(e.TopicChannel, ctx)
	go e.publisher.Run(ctx)
	go e.MsgProcessor.Run(ctx)
	go e.processReadings(ctx)
}

func (e *MimirEngine) processReadings(ctx context.Context) {
	for {
		select {
		case reading := <-e.ReadingChannel:
			go func() {
				db.StoreReading(reading)
			}()
		case <-ctx.Done():
			slog.Info("context done", "error", ctx.Err())
			return
		}
	}
}

func (e *MimirEngine) Close() {
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

func (e *MimirEngine) BuildTrigger(trigger models.TriggerOptions) (triggers.Trigger, error) {
	return e.TriggerFactory.BuildTrigger(trigger)
}
