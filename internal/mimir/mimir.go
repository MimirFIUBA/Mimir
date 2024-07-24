package mimir

import (
	"fmt"
	"time"
)

var (
	// sensorManager = SensorManager{}
	Data = DataManager{}
)

func Run() {

	group := NewGroup("group 1")
	Data.AddGroup(group)

	// Load Config
	//Test Sensor
	sensor := NewSensor("test sensor 1")

	value1 := 1
	data := SensorReading{Value: value1, Time: time.Now()}
	sensor.Data = append(sensor.Data, data)

	value2 := 1.3
	data2 := SensorReading{Value: value2, Time: time.Now()}
	sensor.Data = append(sensor.Data, data2)

	value3 := true
	data3 := SensorReading{Value: value3, Time: time.Now()}
	sensor.Data = append(sensor.Data, data3)

	valueString := "too high"
	data4 := SensorReading{Value: valueString, Time: time.Now()}
	sensor.Data = append(sensor.Data, data4)

	fmt.Printf("sensorData: %+v\n", len(sensor.Data))
	Data.AddSensor(sensor)
	//End test sensor

	Data.AddSensor(NewSensor("sensorPH"))
	// CreateSensor(Sensor{0, "sensorTemp", "this is a temp sensor", nil})
}
