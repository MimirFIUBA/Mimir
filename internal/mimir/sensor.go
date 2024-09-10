package mimir

import (
	"mimir/internal/triggers"
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
	triggerList []triggers.TriggerObserver
}

type SensorReading struct {
	SensorID string      `json:"sensorId"`
	Value    SensorValue `json:"value"`
	Time     time.Time   `json:"time"`
}

type SensorValue interface{}

func NewSensor(name string) *Sensor {
	return &Sensor{"", name, "", "", "", "", nil, nil}
}

func (s *Sensor) addReading(reading SensorReading) {
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

func (s *Sensor) Register(observer triggers.TriggerObserver) {
	s.triggerList = append(s.triggerList, observer)
}

// func (s *Sensor) deregister(observer trigger.Observer) {
// 	s.observerList = removeFromslice(s.observerList, observer)
// }

func (s *Sensor) notifyAll() {
	for _, observer := range s.triggerList {
		event := triggers.Event{
			Name:      "new reading event",
			Timestamp: time.Now(),
			Data:      s.Data[len(s.Data)-1].Value}
		observer.Update(event) //TODO: need to send the last value
	}
}

// func removeFromslice(observerList []trigger.Observer, observerToRemove trigger.Observer) []trigger.Observer {
// 	observerListLength := len(observerList)
// 	for i, observer := range observerList {
// 		if observerToRemove.GetID() == observer.GetID() {
// 			observerList[observerListLength-1], observerList[i] = observerList[i], observerList[observerListLength-1]
// 			return observerList[:observerListLength-1]
// 		}
// 	}
// 	return observerList
// }
