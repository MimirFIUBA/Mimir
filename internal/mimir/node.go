package mimir

type Node struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	GroupID     string   `json:"groupId"`
	Sensors     []Sensor `json:"sensors"`
}

func NewNode(name string) *Node {
	return &Node{"", name, "", "", nil}
}

func (n *Node) Update(updatedNode *Node) {
	n.Name = updatedNode.Name
	n.Description = updatedNode.Description
	n.GroupID = updatedNode.GroupID
}
