package main

import (
	"fmt"
	API "mimir/internal/api"
	mimir "mimir/internal/mimir"
	"mimir/internal/triggers"

	// mqtt "mimir/internal/mqtt"
	"os"
	"os/signal"
	"syscall"
)

func setInitialData(mp *mimir.MimirProcessor) {
	node := mimir.NewNode("esp32")
	node = mimir.Data.AddNode(node)

	// testSensor := NewSensor("test")
	// testSensor.DataName = "mimirTest"
	// Data.AddSensor(testSensor)

	// // processor := &JsonProcessor{"id", "value"}

	// processor := NewBytesProcessor()
	// idConfiguration := NewBytesConfiguration("id", binary.BigEndian, 1)
	// dataConfiguration := NewBytesConfiguration("bool", binary.BigEndian, 1)
	// processor.AddBytesConfiguration(*idConfiguration)
	// processor.AddBytesConfiguration(*dataConfiguration)
	// processor.AddBytesConfiguration(*NewBytesConfiguration("id", binary.BigEndian, 1))
	// processor.AddBytesConfiguration(*NewBytesConfiguration("float", binary.BigEndian, 4))

	// MessageProcessors.RegisterProcessor("mimir/mimirTest", processor)

	// dhtTemperatureSensor := NewSensor("dht temp")
	// dhtTemperatureSensor.DataName = "temperature"
	// dhtTemperatureSensor.NodeID = node.ID

	// processorTemp := NewJSONProcessor()
	// configuration1 := JSONValueConfiguration{"id", "value"}
	// processorTemp.jsonValueConfigurations = append(processorTemp.jsonValueConfigurations, configuration1)
	// MessageProcessors.RegisterProcessor("mimir/esp32/temperature", processorTemp)

	// dhtHumiditySensor := NewSensor("dht humidity")
	// dhtHumiditySensor.DataName = "humidity"
	// dhtHumiditySensor.NodeID = node.ID

	// // Trigger construction
	// printAction := triggers.PrintAction{Message: "TRIGGER EXECUTED - Temperature too low"}
	// sendMQTTMessageAction := triggers.SendMQTTMessageAction{
	// 	Message:                 "Temperature too low on sensor 1",
	// 	OutgoingMessagesChannel: outgoingMessagesChannel}
	// condition := triggers.GenericCondition{Operator: "<", ReferenceValue: 10.0, TestValue: 0.0}

	// sensorTestTrigger := triggers.NewTrigger("sensor test trigger")
	// sensorTestTrigger.Condition = &condition
	// sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &printAction)
	// sensorTestTrigger.Actions = append(sensorTestTrigger.Actions, &sendMQTTMessageAction)

	// dhtTemperatureSensor.register(sensorTestTrigger)

	// printActionTimeTrigger := triggers.PrintAction{Message: "ALERT!!! Time Trigger executed"}
	// sendWSAction := triggers.SendMQTTMessageAction{
	// 	Message:                 "on",
	// 	OutgoingMessagesChannel: wsMsgChan}
	// receiveValueCondition := triggers.ReceiveValueCondition{}
	// sensorTimeTrigger := triggers.NewTimeTrigger("sensor time trigger", 10*time.Second)
	// sensorTimeTrigger.Condition = &receiveValueCondition
	// sensorTimeTrigger.Actions = append(sensorTimeTrigger.Actions, &printActionTimeTrigger)
	// sensorTimeTrigger.Actions = append(sensorTimeTrigger.Actions, &sendWSAction)

	// dhtTemperatureSensor.register(sensorTimeTrigger)
	// sensorTimeTrigger.Start()

	// Data.AddSensor(dhtTemperatureSensor)
	// Data.AddSensor(dhtHumiditySensor)

	// Real use case:
	// CREATE SENSOR
	waterTemperatureSensor := mimir.NewSensor("water temperature")
	waterTemperatureSensor.DataName = "waterTemp"
	waterTemperatureSensor.NodeID = node.ID
	// waterTemperatureSensor = Data.AddSensor(waterTemperatureSensor)

	// MESSAGE PROCESSOR
	processorWaterTemp := mimir.NewJSONProcessor()
	processorWaterTemp.SensorId = "0"
	waterTempConfiguration := mimir.NewJSONValueConfiguration("", "data")
	processorWaterTemp.AddValueConfiguration(waterTempConfiguration)
	mimir.MessageProcessors.RegisterProcessor("mimir/esp32/waterTemp", processorWaterTemp)

	//TRIGGERS
	wtPrintAction := triggers.PrintAction{Message: "TRIGGER EXECUTED - Water temperature too high"}
	wtSendMQTTMessageAction := mp.NewSendMQTTMessageAction("on")
	wtHighTemperatureCondition := triggers.CompareCondition{Operator: ">", ReferenceValue: 50.0, TestValue: 0.0}

	wtTrigger := triggers.NewTrigger("sensor test trigger")
	wtTrigger.Condition = &wtHighTemperatureCondition
	wtTrigger.Actions = append(wtTrigger.Actions, &wtPrintAction)
	wtTrigger.Actions = append(wtTrigger.Actions, wtSendMQTTMessageAction)

	waterTemperatureSensor.Register(wtTrigger)

	wtLowPrintAction := triggers.PrintAction{Message: "TRIGGER EXECUTED - Water temperature low"}
	wtLowSendMQTTMessageAction := mp.NewSendMQTTMessageAction("off")
	wtLowTemperatureCondition := triggers.CompareCondition{Operator: "<", ReferenceValue: 50.0, TestValue: 0.0}

	wtLowTrigger := triggers.NewTrigger("sensor test trigger")
	wtLowTrigger.Condition = &wtLowTemperatureCondition
	wtLowTrigger.Actions = append(wtLowTrigger.Actions, &wtLowPrintAction)
	wtLowTrigger.Actions = append(wtLowTrigger.Actions, wtLowSendMQTTMessageAction)

	waterTemperatureSensor.Register(wtLowTrigger)
	mimir.Data.AddSensor(waterTemperatureSensor)

}

func main() {

	topics := mimir.GetTopics()
	client := mimir.StartMqttClient()

	mimirProcessor := mimir.NewMimirProcessor()

	mimirProcessor.StartGateway(client, topics)
	setInitialData(mimirProcessor)
	go mimirProcessor.Run()
	go API.Start(mimirProcessor.WsChannel)
	go API.StartWebSocket()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	mimir.CloseConnection(client, topics)

	fmt.Println("Mimir is out of duty, bye!")
}
