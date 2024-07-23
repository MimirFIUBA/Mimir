package mimir

import "fmt"

type SensorManager struct {
	sensors []Sensor
}

func (s *SensorManager) getNewSensorId() int {
	return len(s.sensors)
}

func (s *SensorManager) storeReading(reading SensorReading) {
	for i := range s.sensors {
		sensor := &s.sensors[i]
		if sensor.ID == reading.SensorID {
			sensor.addReading(reading)
			break
		}
	}
}

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
