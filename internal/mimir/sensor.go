package mimir

import (
	"time"
)

type Sensor struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	DataName    string          `json:"dataName"`
	Topic       string          `json:"topic"`
	NodeID      string          `json:"nodeId"`
	Description string          `json:"description"`
	Data        []SensorReading `json:"data"`
}

type SensorReading struct {
	SensorID string      `json:"sensorId"`
	Value    SensorValue `json:"value"`
	Time     time.Time   `json:"time"`
}

type SensorValue interface{}

func NewSensor(name string) *Sensor {
	return &Sensor{"", name, "", "", "", "", nil}
}

func (s *Sensor) addReading(reading SensorReading) {
	s.Data = append(s.Data, reading)
}

func (s *Sensor) Update(newData *Sensor) {
	//TODO: check the best way to do this.
	s.Name = newData.Name
	s.DataName = newData.DataName
	s.Topic = newData.Topic
	s.NodeID = newData.NodeID
	s.Description = newData.Description
}

func (s *Sensor) GetId() string {
	return s.ID
}
