package mimir

import (
	"fmt"
)

var (
	SensorChannel = make(chan Sensor, 1)
	sensorManager = SensorManager{}
)

type Sensor struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}

type SensorManager struct {
	sensors []Sensor
}

func (s *SensorManager) getNewSensorId() int {
	return len(s.sensors)
}

func CreateSensor(sensor Sensor) Sensor {
	sensor.ID = sensorManager.getNewSensorId()
	sensorManager.sensors = append(sensorManager.sensors, sensor)
	fmt.Printf("New sensor created: %+v\n", sensor) 
	return sensor
}

func GetSensors() []Sensor {
	return sensorManager.sensors
}

func Run () {
}