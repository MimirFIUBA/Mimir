package models

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
