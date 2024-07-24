package mimir

import "github.com/google/uuid"

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Nodes       []Node    `json:"nodes"`
}

func NewGroup(name string) *Group {
	id := uuid.New()
	return &Group{id, name, "", "crop", nil}
}
