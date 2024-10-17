package models

import (
	"mimir/triggers"
	"slices"
	"time"
)

type Sensor struct {
	ID          string          `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string          `json:"name" bson:"name"`
	DataName    string          `json:"dataName" bson:"data_name"`
	Topic       string          `json:"topic" bson:"topic"`
	NodeID      string          `json:"nodeId" bson:"node_id"`
	Description string          `json:"description" bson:"description"`
	IsActive    bool            `json:"isActive" bson:"is_active"`
	Data        []SensorReading `json:"data" bson:"data, omitempty"`
	triggerList []triggers.TriggerObserver
}

type SensorReading struct {
	SensorID string      `json:"sensorId"`
	Topic    string      `json:"topic"`
	Value    SensorValue `json:"value"`
	Time     time.Time   `json:"time"`
}

type SensorValue interface{}

func NewSensor(name string) *Sensor {
	return &Sensor{Name: name}
}

func (s *Sensor) AddReading(reading SensorReading) {
	s.Data = append(s.Data, reading)
	s.notifyAll()
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

func (s *Sensor) Register(trigger triggers.TriggerObserver) {
	s.triggerList = append(s.triggerList, trigger)
}

func (s *Sensor) Deregister(trigger triggers.TriggerObserver) {
	idToRemove := trigger.GetID()

	s.triggerList = slices.DeleteFunc(s.triggerList, func(trigger triggers.TriggerObserver) bool {
		return trigger.GetID() == idToRemove
	})
}

func (s *Sensor) notifyAll() {
	for _, observer := range s.triggerList {
		reading := s.Data[len(s.Data)-1]
		event := triggers.Event{
			Name:      "new reading event",
			Timestamp: time.Now(),
			Data:      reading.Value,
			SenderId:  reading.Topic,
		}
		observer.Update(event) //TODO: need to send the last value
	}
}

func (s *Sensor) GetTriggers() []triggers.TriggerObserver {
	return s.triggerList
}
