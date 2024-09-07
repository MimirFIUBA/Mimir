package mimir

import (
	"encoding/binary"
	"fmt"
	"mimir/internal/triggers"
	"time"
)

var (
	Data = DataManager{}
)

type MimirProcessor struct {
	outgoingMessagesChannel chan string
	readingChannel          chan SensorReading
	topicChannel            chan string
	wsChannel               chan string
}

func NewMimirProcessor(topicChannel chan string, readingChannel chan SensorReading, outgoingMessagesChannel chan string, wsChannel chan string) *MimirProcessor {
	Data.topicChannel = topicChannel
	mp := MimirProcessor{outgoingMessagesChannel, readingChannel, topicChannel, wsChannel}
	return &mp
}

func (mp *MimirProcessor) Run() {
	setInitialData(mp.outgoingMessagesChannel, mp.wsChannel)

	for {
		reading := <-mp.readingChannel

		processReading(reading)

		Data.StoreReading(reading)
	}
}

func processReading(reading SensorReading) {
	//TODO: ver si necesitamos hacer algo aca
	fmt.Printf("Processing reading: %v \n", reading.Value)
}

func setInitialData(outgoingMessagesChannel chan string, wsMsgChan chan string) {
	testSensor := NewSensor("test")
	testSensor.DataName = "mimirTest"
	Data.AddSensor(testSensor)

	// processor := &JsonProcessor{"id", "value"}

	processor := NewBytesProcessor()
	idConfiguration := NewBytesConfiguration("id", binary.BigEndian, 1)
	dataConfiguration := NewBytesConfiguration("bool", binary.BigEndian, 1)
	processor.AddBytesConfiguration(*idConfiguration)
	processor.AddBytesConfiguration(*dataConfiguration)
	processor.AddBytesConfiguration(*NewBytesConfiguration("id", binary.BigEndian, 1))
	processor.AddBytesConfiguration(*NewBytesConfiguration("float", binary.BigEndian, 4))

	MessageProcessors.RegisterProcessor("mimir/mimirTest", processor)

	node := NewNode("esp32")
	node = Data.AddNode(node)

	dhtTemperatureSensor := NewSensor("dht temp")
	dhtTemperatureSensor.DataName = "temperature"
	dhtTemperatureSensor.NodeID = node.ID

	processorTemp := NewJSONProcessor()
	configuration1 := JSONValueConfiguration{"id", "value"}
	processorTemp.jsonValueConfigurations = append(processorTemp.jsonValueConfigurations, configuration1)
	MessageProcessors.RegisterProcessor("mimir/esp32/temperature", processorTemp)

	dhtHumiditySensor := NewSensor("dht humidity")
	dhtHumiditySensor.DataName = "humidity"
	dhtHumiditySensor.NodeID = node.ID

	// Trigger construction
	printAction := triggers.PrintAction{Message: "TRIGGER EXECUTED - Temperature too low"}
	sendMQTTMessageAction := triggers.SendMQTTMessageAction{
		Message:                 "Temperature too low on sensor 1",
		OutgoingMessagesChannel: outgoingMessagesChannel}
	condition := triggers.GenericCondition{Operator: "<", ReferenceValue: 10.0, TestValue: 0.0}

	sensorTestTrigger := triggers.NewTrigger("sensor test trigger")
	sensorTestTrigger.Condition = &condition
	sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &printAction)
	sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &sendMQTTMessageAction)

	dhtTemperatureSensor.register(sensorTestTrigger)

	printActionTimeTrigger := triggers.PrintAction{Message: "ALERT!!! Time Trigger executed"}
	sendWSAction := triggers.SendMQTTMessageAction{
		Message:                 "{Hello}",
		OutgoingMessagesChannel: wsMsgChan}
	receiveValueCondition := triggers.ReceiveValueCondition{}
	sensorTimeTrigger := triggers.NewTimeTrigger("sensor time trigger", 10*time.Second)
	sensorTimeTrigger.Condition = &receiveValueCondition
	sensorTimeTrigger.Actions = append(sensorTimeTrigger.Actions, &printActionTimeTrigger)
	sensorTimeTrigger.Actions = append(sensorTimeTrigger.Actions, &sendWSAction)

	dhtTemperatureSensor.register(sensorTimeTrigger)
	sensorTimeTrigger.Start()

	Data.AddSensor(dhtTemperatureSensor)
	Data.AddSensor(dhtHumiditySensor)
}
