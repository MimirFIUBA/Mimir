package models

type Sensor struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DataName    string `json:"dataName"`
	Topic       string `json:"topic"`
	NodeID      string `json:"nodeId"`
	Description string `json:"description"`
}

func (s *Sensor) Update(newData *Sensor) {
	s.Name = newData.Name
	s.DataName = newData.DataName
	s.Topic = newData.Topic
	s.NodeID = newData.NodeID
	s.Description = newData.Description
}
