package responses

import (
	"mimir/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorResponse struct {
	ID                string                `json:"id"`
	Name              string                `json:"name"`
	DataName          string                `json:"dataName"`
	Topic             string                `json:"topic"`
	NodeID            string                `json:"nodeId"`
	Node              NodeResponse          `json:"node"`
	Description       string                `json:"description"`
	IsActive          bool                  `json:"isActive"`
	LastSensedReading *models.SensorReading `json:"lastSensedReading"`
}

func NewSensorResponse(sensor models.Sensor) *SensorResponse {
	return &SensorResponse{
		ID:                sensor.ID.Hex(),
		Name:              sensor.Name,
		DataName:          sensor.DataName,
		Topic:             sensor.Topic,
		NodeID:            sensor.NodeID,
		Description:       sensor.Description,
		IsActive:          sensor.IsActive,
		LastSensedReading: sensor.LastSensedReading,
	}
}

type NodeResponse struct {
	ID          string        `json:"id,omitempty"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	GroupID     string        `json:"groupId"`
	Group       GroupResponse `json:"group"`
}

func NewNodeResponse(node models.Node) *NodeResponse {
	return &NodeResponse{
		ID:          node.ID.Hex(),
		Name:        node.Name,
		Description: node.Description,
		GroupID:     node.GroupID,
	}
}

type GroupResponse struct {
	ID          primitive.ObjectID `json:"id,omitempty"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Type        string             `json:"type"`
}

func NewGroupResponse(group models.Group) *GroupResponse {
	return &GroupResponse{
		ID:          group.ID,
		Name:        group.Name,
		Description: group.Description,
		Type:        group.Type,
	}
}
