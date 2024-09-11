package db

import (
	"fmt"
	"mimir/internal/api/models"
	"strconv"
)

type NodesManager struct {
	nodes     []models.Node
	idCounter int
}

func (n *NodesManager) GetNewId() int {
	n.idCounter++
	return n.idCounter
}

func (n *NodesManager) GetNodes() []models.Node {
	return n.nodes
}

func (n *NodesManager) GetNodeById(id string) (*models.Node, error) {
	for index, node := range n.nodes {
		if node.ID == id {
			return &n.nodes[index], nil
		}
	}

	// TODO(#19) - Improve error handling
	return nil, fmt.Errorf("node %s not found", id)
}

func (n *NodesManager) CreateNode(node *models.Node) error {
	newId := n.GetNewId()
	node.ID = strconv.Itoa(newId)

	n.nodes = append(n.nodes, *node)
	return nil
}

func (n *NodesManager) UpdateNode(node *models.Node) (*models.Node, error) {
	// TODO(#19) - Improve error handling
	oldNode, err := n.GetNodeById(node.ID)
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
		if node.ID == id {
			nodeIndex = i
			break
		}
	}

	// TODO(#19) - Improve error handling
	if nodeIndex == -1 {
		return fmt.Errorf("sensor %s not found", id)
	}

	n.nodes[nodeIndex] = n.nodes[len(n.nodes)-1]
	n.nodes = n.nodes[:len(n.nodes)-1]
	return nil
}

func (n *NodesManager) AddSensorToNodeById(id string, sensor *models.Sensor) error {
	// TODO(#19) - Improve error handling
	oldNode, err := n.GetNodeById(id)
	if err != nil {
		return nil
	}

	return oldNode.AddSensor(sensor)
}
