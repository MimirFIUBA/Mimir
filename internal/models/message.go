package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MessageType string

const (
	ALERT_MESSAGE_TYPE MessageType = "ALERT"
	INFO_MESSAGE_TYPE  MessageType = "INFO"
)

type Message struct {
	Id                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Body              string             `json:"body" bson:"body,omitempty"`
	Type              MessageType        `json:"type" bson:"type,omitempty"`
	IsRead            bool               `json:"read" bson:"read,omitempty"`
	CreatedDate       time.Time          `json:"createdDate" bson:"createdDate,omitempty"`
	AdditionalDetails interface{}        `json:"additionalDetails,omitempty" bson:"additionalDetails,omitempty"`
}

type MqttOutgoingMessage struct {
	Topic   string
	Message string
}

func NewMqttOutgoingMessage(topic, message string) *MqttOutgoingMessage {
	return &MqttOutgoingMessage{topic, message}
}

type WSOutgoingMessage struct {
	Type    string
	Message string
}
