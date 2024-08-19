package mimir

import "fmt"

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
	sensor := NewSensor("sensorPH")
	sensor.DataName = "ph"

	condition := MaxValueCondition{10.0, nil}
	printAction := PrintAction{}
	sendMQTTMessageAction := SendMQTTMessageAction{"Alert test!", outgoingMessagesChannel}

	actions := []Action{&printAction, &sendMQTTMessageAction}
	trigger := Trigger{&condition, actions}
	sensor.Triggers = append(sensor.Triggers, trigger)

	Data.AddSensor(sensor)
}

func (mp *MimirProcessor) Run() {
	setInitialData(mp.outgoingMessagesChannel)

	for {
		reading := <-mp.readingChannel

		//TODO: processReading
		processReading(reading)

		Data.StoreReading(reading)
	}
}

func processReading(reading SensorReading) {
	fmt.Printf("Processing reading: %v \n", reading.Value)
	sensor := Data.GetSensor(reading.SensorID)
	for _, trigger := range sensor.Triggers {
		trigger.Execute(reading)
	}
}

// func Run(topicChannel chan string) {
// 	Data.topicChannel = topicChannel

// 	group := NewGroup("group 1")
// 	Data.AddGroup(group)

// 	// Load Config
// 	//Test Sensor
// 	sensor := NewSensor("test sensor 1")

// 	value1 := 1
// 	data := SensorReading{Value: value1, Time: time.Now()}
// 	sensor.Data = append(sensor.Data, data)

// 	value2 := 1.3
// 	data2 := SensorReading{Value: value2, Time: time.Now()}
// 	sensor.Data = append(sensor.Data, data2)

// 	value3 := true
// 	data3 := SensorReading{Value: value3, Time: time.Now()}
// 	sensor.Data = append(sensor.Data, data3)

// 	valueString := "too high"
// 	data4 := SensorReading{Value: valueString, Time: time.Now()}
// 	sensor.Data = append(sensor.Data, data4)

// 	fmt.Printf("sensorData: %+v\n", len(sensor.Data))
// 	Data.AddSensor(sensor)
// 	//End test sensor

// 	Data.AddSensor(NewSensor("sensorPH"))
// 	// CreateSensor(Sensor{0, "sensorTemp", "this is a temp sensor", nil})

// }
