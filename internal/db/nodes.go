package db

import (
	"context"
	"fmt"
	"log/slog"
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (n *NodesManager) CreateNode(node *models.Node) error {
	// TODO(#20) - Add Body validation
	node, err := Database.insertNode(node)
	if err != nil {
		slog.Error("error inserting node", "error", err)
		return err
	}

	n.nodes = append(n.nodes, *node)
	err = GroupsData.AddNodeToGroupById(node.GroupID, node)
	if err != nil {
		return err
	}

	return nil
}

func (n *NodesManager) UpdateNode(node *models.Node) (*models.Node, error) {
	oldNode, err := n.GetNodeById(node.GetId())
	if err != nil {
		return nil, err
	}

	oldNode.Update(node)

	_, err = Database.UpdateNode(oldNode)
	if err != nil {
		return nil, err
	}

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

func (n *NodesManager) AddSensorToNodeById(id string, sensor *models.Sensor) error {
	oldNode, err := n.GetNodeById(id)
	if err != nil {
		return nil
	}

	err = oldNode.AddSensor(sensor)

	if err != nil {
		return err
	}

	_, err = Database.UpdateNode(oldNode)

	return err
}

func (n *NodesManager) AddNode(node *models.Node) {
	n.nodes = append(n.nodes, *node)
}

func (d *DatabaseManager) insertNode(node *models.Node) (*models.Node, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		nodesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(NODES_COLLECTION)
		result, err := nodesCollection.InsertOne(context.TODO(), node)
		if err != nil {
			return nil, err
		}

		nodeId, ok := result.InsertedID.(primitive.ObjectID)
		if !ok {
			return nil, fmt.Errorf("error converting id for node")
		}
		node.ID = nodeId
	}

	return node, nil
}

func (d *DatabaseManager) UpdateNode(node *models.Node) (*models.Node, error) {
	mongoClient := d.getMongoClient()
	if mongoClient != nil {
		nodesCollection := mongoClient.Database(MONGO_DB_MIMIR).Collection(NODES_COLLECTION)
		objectId, err := primitive.ObjectIDFromHex(node.GetId())
		if err != nil {
			return nil, err
		}
		update := bson.D{{Key: "$set", Value: node}}
		_, err = nodesCollection.UpdateByID(context.TODO(), objectId, update)
		if err != nil {
			fmt.Println("error updating node ", err)
			return nil, err
		}
	}
	return node, nil
}
