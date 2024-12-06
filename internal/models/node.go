package models

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Node struct {
	ID          primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	GroupID     string             `json:"groupId" bson:"group_id"`
	Sensors     []*Sensor          `json:"sensors,omitempty" bson:"sensors, omitempty"`
}

func NewNode(name string) *Node {
	return &Node{primitive.ObjectID{}, name, "", "", nil}
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

	n.Sensors = append(n.Sensors, sensor)

	return nil
}

func (n *Node) GetId() string {
	return n.ID.Hex()
}
