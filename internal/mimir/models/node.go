package models

import "fmt"

type Node struct {
	ID          string   `json:"id" bson:"mimir_id, omitempty"`
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	GroupID     string   `json:"groupId" bson:"group_id"`
	Sensors     []Sensor `json:"sensors" bson:"sensors, omitempty"`
}

func NewNode(name string) *Node {
	return &Node{"", name, "", "", nil}
}

func (n *Node) Update(updatedNode *Node) {
	n.Name = updatedNode.Name
	n.Description = updatedNode.Description
	n.GroupID = updatedNode.GroupID
}

func (n *Node) AddSensor(sensor *Sensor) error {
	for _, s := range n.Sensors {
		if s.ID == sensor.ID {
			return fmt.Errorf("already exists sensor")
		}
	}

	n.Sensors = append(n.Sensors, *sensor)

	sensor.Topic = "mimir/" + n.Name + "/" + sensor.DataName

	return nil
}
