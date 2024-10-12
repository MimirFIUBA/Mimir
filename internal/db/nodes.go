package db

import (
	"fmt"
	mimir "mimir/internal/mimir/models"
)

type NodesManager struct {
	nodes     []mimir.Node
	idCounter int
}

func (n *NodesManager) GetNewId() int {
	n.idCounter++
	return n.idCounter
}

func (n *NodesManager) GetNodes() []mimir.Node {
	return n.nodes
}

func (n *NodesManager) GetNodeById(id string) (*mimir.Node, error) {
	for index, node := range n.nodes {
		if node.GetId() == id {
			return &n.nodes[index], nil
		}
	}

	return nil, fmt.Errorf("node %s not found", id)
}

func (n *NodesManager) IdExists(id string) bool {
	_, err := n.GetNodeById(id)
	return err == nil
}

func (n *NodesManager) CreateNode(node *mimir.Node) error {
	// TODO(#20) - Add Body validation
	node, err := Database.insertNode(node)
	if err != nil {
		return err
	}

	n.nodes = append(n.nodes, *node)
	err = GroupsData.AddNodeToGroupById(node.GroupID, node)
	if err != nil {
		return err
	}

	return nil
}

func (n *NodesManager) UpdateNode(node *mimir.Node) (*mimir.Node, error) {
	oldNode, err := n.GetNodeById(node.GetId())
	if err != nil {
		return nil, err
	}

	oldNode.Update(node)
	return oldNode, nil
}

func (n *NodesManager) DeleteNode(id string) error {
	nodeIndex := -1
	for i := range n.nodes {
		node := &n.nodes[i]
		if node.GetId() == id {
			nodeIndex = i
			break
		}
	}

	if nodeIndex == -1 {
		return fmt.Errorf("sensor %s not found", id)
	}

	n.nodes[nodeIndex] = n.nodes[len(n.nodes)-1]
	n.nodes = n.nodes[:len(n.nodes)-1]
	return nil
}

func (n *NodesManager) AddSensorToNodeById(id string, sensor *mimir.Sensor) error {
	oldNode, err := n.GetNodeById(id)
	if err != nil {
		return nil
	}

	return oldNode.AddSensor(sensor)
}

func (n *NodesManager) AddNode(node *mimir.Node) {
	n.nodes = append(n.nodes, *node)
}
