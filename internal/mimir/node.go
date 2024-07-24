package mimir

import "github.com/google/uuid"

type Node struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	GroupID     uuid.UUID `json:"groupId"`
	Sensors     []Sensor  `json:"data"`
}

func NewNode(name string) *Node {
	id := uuid.New()
	return &Node{id, name, "", uuid.Nil, nil}
}
