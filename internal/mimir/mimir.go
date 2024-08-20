package mimir

import (
	"fmt"
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

func setInitialData(outgoingMessagesChannel chan string) {
	sensor := NewSensor("sensorTest")
	sensor.DataName = "mimirTest"

	// Test Data
	condition := GenericCondition{10.0, nil, "<"}
	printAction := PrintAction{"Action executed"}
	sendMQTTMessageAction := SendMQTTMessageAction{"Alert test!", outgoingMessagesChannel}

	actions := []Action{&printAction, &sendMQTTMessageAction}
	trigger := Trigger{&condition, actions}
	sensor.Triggers = append(sensor.Triggers, trigger)

	receiveValueCondition := ReceiveValueCondition{}
	timeTriggerActions := []Action{&printAction}

	timeTrigger := NewTimeTrigger(&receiveValueCondition, timeTriggerActions, 5*time.Second)
	timeTrigger.Start()

	sensor.TimeTriggers = append(sensor.TimeTriggers, *timeTrigger)

	Data.AddSensor(sensor)
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
	sensor := Data.GetSensor(reading.SensorID)
	for _, timeTrigger := range sensor.TimeTriggers {
		timeTrigger.Evaluate(reading)
	}

	for _, trigger := range sensor.Triggers {
		trigger.Execute(reading)
	}

}
