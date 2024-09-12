package mimir

type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Nodes       []Node `json:"nodes"`
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
