package mimir

import (
	"fmt"
	"mimir/internal/trigger"
	"time"
)

var (
	Data = DataManager{}
)

type MimirProcessor struct {
	outgoingMessagesChannel chan string
	readingChannel          chan SensorReading
	topicChannel            chan string
}

func NewMimirProcessor(topicChannel chan string, readingChannel chan SensorReading, outgoingMessagesChannel chan string) *MimirProcessor {
	Data.topicChannel = topicChannel
	mp := MimirProcessor{outgoingMessagesChannel, readingChannel, topicChannel}
	return &mp
}

func (mp *MimirProcessor) Run() {
	setInitialData(mp.outgoingMessagesChannel)

	for {
		reading := <-mp.readingChannel

		processReading(reading)

		Data.StoreReading(reading)
	}
}

func processReading(reading SensorReading) {
	fmt.Printf("Processing reading: %v \n", reading.Value)
	// sensor := Data.GetSensor(reading.SensorID)
	// for _, timeTrigger := range sensor.TimeTriggers {
	// 	timeTrigger.Evaluate(reading)
	// }

	// for _, trigger := range sensor.Triggers {
	// 	trigger.Execute(reading)
	// }
}

func setInitialData(outgoingMessagesChannel chan string) {
	sensor := NewSensor("test")
	sensor.DataName = "mimirTest"

	// Trigger construction
	printAction := trigger.PrintAction{Message: "Action executed"}
	sendMQTTMessageAction := trigger.SendMQTTMessageAction{
		Message:                 "Temperature too low on sensor 1",
		OutgoingMessagesChannel: outgoingMessagesChannel}
	condition := trigger.GenericCondition{Operator: "<", ReferenceValue: 10.0, TestValue: 0.0}

	sensorTestTrigger := trigger.NewTrigger("sensor test trigger")
	sensorTestTrigger.Condition = &condition
	sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &printAction)
	sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &sendMQTTMessageAction)

	sensor.register(sensorTestTrigger)

	printActionTimeTrigger := trigger.PrintAction{Message: "Time Trigger executed"}
	receiveValueCondition := trigger.ReceiveValueCondition{}
	sensorTimeTrigger := trigger.NewTimeTrigger("sensor time trigger", 10*time.Second)
	sensorTimeTrigger.Condition = &receiveValueCondition
	sensorTimeTrigger.Actions = append(sensorTimeTrigger.Actions, &printActionTimeTrigger)

	sensor.register(sensorTimeTrigger)
	sensorTimeTrigger.Start()

	Data.AddSensor(sensor)
}
