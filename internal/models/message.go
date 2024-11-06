package models

import "time"

type MessageType string

const (
	ALERT_MESSAGE_TYPE MessageType = "ALERT"
	INFO_MESSAGE_TYPE  MessageType = "INFO"
)

type Message struct {
	Id                string      `json:"id,omitempty" bson:"_id,omitempty"`
	Body              string      `json:"body" bson:"body"`
	Type              MessageType `json:"type" bson:"type"`
	IsRead            bool        `json:"read" bson:"read"`
	CreatedDate       time.Time   `json:"createdDate" bson:"createdDate"`
	AdditionalDetails interface{} `json:"additionalDetails,omitempty" bson:"additionalDetails,omitempty"`
}

type MqttOutgoingMessage struct {
	Topic   string
	Message string
}

func NewMqttOutgoingMessage(topic, message string) *MqttOutgoingMessage {
	return &MqttOutgoingMessage{topic, message}
}
