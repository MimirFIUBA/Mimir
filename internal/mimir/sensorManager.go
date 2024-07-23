package mimir

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
