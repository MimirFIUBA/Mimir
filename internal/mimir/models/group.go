package models

import "errors"

type Group struct {
	ID          string `json:"id" bson:"mimirId, omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Type        string `json:"type" bson:"type"`
	Nodes       []Node `json:"nodes" bson:"nodes, omitempty"`
}

func NewGroup(name string) *Group {
	return &Group{"", name, "", "crop", nil}
}

func (g *Group) Update(updatedGroup *Group) {
	g.Name = updatedGroup.Name
	g.Description = updatedGroup.Description
	g.Type = updatedGroup.Type
}

func (g *Group) GetId() string {
	return g.ID
}

func (g *Group) AddNode(node *Node) error {
	for _, n := range g.Nodes {
		if n.ID == node.ID {
			return errors.New("node already exists")
		}
	}

	g.Nodes = append(g.Nodes, *node)
	return nil
}
