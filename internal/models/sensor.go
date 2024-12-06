package models

import (
	"fmt"
	"log/slog"
	"mimir/triggers"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sensor struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name              string             `json:"name" bson:"name"`
	DataName          string             `json:"dataName" bson:"data_name"`
	DataLabel         string             `json:"dataLabel" bson:"data_label"`
	Unit              string             `json:"unit" bson:"unit"`
	Topic             string             `json:"topic" bson:"topic"`
	NodeID            string             `json:"nodeId" bson:"node_id"`
	Description       string             `json:"description" bson:"description"`
	IsActive          bool               `json:"isActive" bson:"is_active"`
	Data              []SensorReading    `json:"data" bson:"data, omitempty"`
	LastSensedReading SensorReading      `json:"lastSensedReading" bson:"lastSensedReading"`
	triggerList       []triggers.Trigger
}

func NewSensor(name string) *Sensor {
	return &Sensor{Name: name}
}

func (s *Sensor) AddReading(reading SensorReading) {
	s.Data = append(s.Data, reading)
	s.LastSensedReading = reading
	fmt.Println("")
	s.NotifyAll()
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
	return s.ID.Hex()
}

func (s *Sensor) Register(trigger triggers.Trigger) {
	s.triggerList = append(s.triggerList, trigger)
	trigger.AddSubject(s)
}

func (s *Sensor) Deregister(trigger triggers.Trigger) {
	idToRemove := trigger.GetID()

	s.triggerList = slices.DeleteFunc(s.triggerList, func(trigger triggers.Trigger) bool {
		return trigger.GetID() == idToRemove
	})
}

func (s *Sensor) NotifyAll() {
	for _, trigger := range s.triggerList {
		reading := s.LastSensedReading
		data := make(map[string]interface{})
		data["reading"] = reading

		event := triggers.Event{
			Name:      "new reading event",
			Timestamp: time.Now(),
			Data:      data,
			Value:     reading.Value,
			SenderId:  reading.Topic,
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("error with trigger update", "error", r)
				}
			}()
			trigger.Update(event)
		}()
	}
}

func (s *Sensor) GetTriggers() []triggers.Trigger {
	return s.triggerList
}
