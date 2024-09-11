package models

import "fmt"

type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	GroupID     string   `json:"groupId"`
	Sensors     []Sensor `json:"sensors"`
}

func (n *Node) Update(updatedNode *Node) {
	n.Name = updatedNode.Name
	n.Description = updatedNode.Description
	n.GroupID = updatedNode.GroupID
}

func (n *Node) AddSensor(sensor *Sensor) error {
	// TODO(#19) - Improve error handling
	for _, s := range n.Sensors {
		if s.ID == sensor.ID {
			return fmt.Errorf("already exists sensor")
		}
	}

	n.Sensors = append(n.Sensors, *sensor)
	return nil
}