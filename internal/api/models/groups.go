package models

import "errors"

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Nodes       []Node `json:"nodes"`
}

func (g *Group) Update(updatedGroup *Group) {
	g.Name = updatedGroup.Name
	g.Description = updatedGroup.Description
	g.Type = updatedGroup.Type
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
