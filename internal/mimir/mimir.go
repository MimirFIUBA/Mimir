package mimir

import (
	"fmt"
	"time"
)

var (
	sensorManager = SensorManager{}
	dataManager   = DataManager{}
)

func AddSensor(sensor *Sensor) *Sensor {
	sensor.ID = sensorManager.getNewSensorId()
	sensorManager.sensors = append(sensorManager.sensors, *sensor)
	fmt.Printf("New sensor created: %+v\n", sensor)
	return sensor
}

func GetSensors() []Sensor {
	return sensorManager.sensors
}

func GetSensor(id int) *Sensor {
	for _, sensor := range sensorManager.sensors {
		if sensor.ID == id {
			return &sensor
		}
	}
	return nil
}

func StoreReading(reding SensorReading) {
	sensorManager.storeReading(reding)
}

func Run() {
	// Load Config
	//Test Sensor
	sensor := NewSensor("test sensor 1")

	var value1 SensorValue
	value1 = 1
	data := SensorReading{Value: value1, Time: time.Now()}
	sensor.Data = append(sensor.Data, data)

	var value2 SensorValue
	value2 = 1.3
	data2 := SensorReading{Value: value2, Time: time.Now()}
	sensor.Data = append(sensor.Data, data2)

	var value3 SensorValue
	value3 = true
	data3 := SensorReading{Value: value3, Time: time.Now()}
	sensor.Data = append(sensor.Data, data3)

	var valueString SensorValue
	valueString = "too high"
	data4 := SensorReading{Value: valueString, Time: time.Now()}
	sensor.Data = append(sensor.Data, data4)

	fmt.Printf("sensorData: %+v\n", len(sensor.Data))
	AddSensor(sensor)
	//End test sensor

	AddSensor(NewSensor("sensorPH"))
	// CreateSensor(Sensor{0, "sensorTemp", "this is a temp sensor", nil})
}
