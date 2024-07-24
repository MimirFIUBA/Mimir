package mimir

import (
	"time"

	"github.com/google/uuid"
)

type Sensor struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	DataName    string          `json:"dataName"`
	NodeID      uuid.UUID       `json:"nodeId"`
	Description string          `json:"description"`
	Data        []SensorReading `json:"data"`
}

type SensorReading struct {
	SensorID int         `json:"sensorId"`
	Value    SensorValue `json:"value"`
	Time     time.Time   `json:"time"`
}

type SensorValue interface{}

func NewSensor(name string) *Sensor {
	return &Sensor{0, name, "", uuid.Nil, "", nil}
}

func (s *Sensor) addReading(reading SensorReading) {
	s.Data = append(s.Data, reading)
}
