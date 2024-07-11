package mimir

type Sensor struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Data []SensorReading `json:"data"`
}

func (s *Sensor) addReading(reading SensorReading) {
	s.Data = append(s.Data, reading)
}