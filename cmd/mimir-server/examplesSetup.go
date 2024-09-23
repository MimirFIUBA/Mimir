package main

import (
	"fmt"
	"mimir/internal/db"
	"mimir/internal/mimir"
	"mimir/internal/mimir/models"
	"mimir/triggers"
	"time"
)

func Setup1(mp *mimir.MimirProcessor) {
	node := models.NewNode("esp32")
	// node = mimir.Data.AddNode(node)
	db.NodesData.CreateNode(node)

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
	waterTemperatureSensor := models.NewSensor("water temperature")
	waterTemperatureSensor.DataName = "waterTemp"
	waterTemperatureSensor.NodeID = node.ID
	// waterTemperatureSensor = Data.AddSensor(waterTemperatureSensor)

	// MESSAGE PROCESSOR
	processorWaterTemp := mimir.NewJSONProcessor()
	processorWaterTemp.SensorId = "1"
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
	wtTrigger.Actions = append(wtTrigger.Actions, &wtSendMQTTMessageAction)

	waterTemperatureSensor.Register(wtTrigger)

	wtLowPrintAction := triggers.PrintAction{Message: "TRIGGER EXECUTED - Water temperature low"}
	wtLowSendMQTTMessageAction := mp.NewSendMQTTMessageAction("off")
	wtLowTemperatureCondition := triggers.CompareCondition{Operator: "<", ReferenceValue: 50.0, TestValue: 0.0}

	wtLowTrigger := triggers.NewTrigger("sensor test trigger")
	wtLowTrigger.Condition = &wtLowTemperatureCondition
	wtLowTrigger.Actions = append(wtLowTrigger.Actions, &wtLowPrintAction)
	wtLowTrigger.Actions = append(wtLowTrigger.Actions, &wtLowSendMQTTMessageAction)

	waterTemperatureSensor.Register(wtLowTrigger)

	//Send through ws trigger
	wsTrigger := triggers.NewTrigger("send ws")
	sendWSAction := mp.NewSendWebSocketMessageAction("")
	sendWSAction.MessageContructor = func(e triggers.Event) string {
		return fmt.Sprintf("newReading: %v", e)
	}
	wsTrigger.Actions = append(wsTrigger.Actions, &sendWSAction)
	waterTemperatureSensor.Register(wsTrigger)

	// freqTrigger := triggers.NewFrequencyTrigger("freq trigger", 3*time.Second)
	// freqTrigger.Actions = append(freqTrigger.Actions, &wtPrintAction)
	// waterTemperatureSensor.Register(freqTrigger)

	// timeoutTrigger := triggers.NewTimeTrigger("tt", 1*time.Second)
	// timeoutTrigger.Condition = &triggers.ReceiveValueCondition{}

	// cmdAction := &triggers.CommandAction{Command: "ls", CommandArgs: "-ltr"}
	// timeoutTrigger.Actions = append(timeoutTrigger.Actions, cmdAction)
	// waterTemperatureSensor.Register(timeoutTrigger)
	// timeoutTrigger.Start()

	db.SensorsData.CreateSensor(waterTemperatureSensor)
	mp.RegisterSensor(waterTemperatureSensor)
}

func testAverageTriggerCondition() {
	avgCondition := triggers.NewAverageCondition(2, 5, 5*time.Second)
	avgCondition.Condition = triggers.NewCompareCondition(">", 5.0)
	trigger := triggers.NewTrigger("testAvg")
	printAction := &triggers.PrintAction{}
	printAction.SetMessageConstructor(func(e triggers.Event) string { return fmt.Sprintf("event: %v", e) })
	trigger.AddAction(printAction)
	trigger.SetCondition(avgCondition)

	for i := range 10 {
		fmt.Println(i)
		event := triggers.Event{Name: "", Timestamp: time.Now(), Data: float64(i)}
		trigger.Update(event)
		if i == 6 {
			time.Sleep(6 * time.Second)
		}
	}
}
