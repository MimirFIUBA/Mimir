package mimir

type Group struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Nodes       []Node `json:"data"`
}

func NewGroup(name string) *Group {
	return &Group{0, name, "", "crop", nil}
}
