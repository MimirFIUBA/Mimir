package factories

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"mimir/internal/models"
	"mimir/triggers"
	"regexp"
	"strings"
)

func newMqttMessageBuilder(topic, message string) func(triggers.Event) models.MqttOutgoingMessage {
	return func(event triggers.Event) models.MqttOutgoingMessage {
		var buffer bytes.Buffer
		re := regexp.MustCompile(`{{(.*?)}}`)

		lastIndex := 0
		for _, match := range re.FindAllStringSubmatchIndex(message, -1) {
			buffer.WriteString(message[lastIndex:match[0]])

			variableName := message[match[2]:match[3]]

			switch {
			case strings.HasPrefix(variableName, "$userVariable."):
				variableValue := getUserVariable(variableName)
				buffer.WriteString(variableValue)
			case strings.HasPrefix(variableName, "$event"):
				variableValue, err := getEventVariable(variableName, event)
				if err != nil {
					slog.Error("error writing value from event to message", "error", err, "variable name", variableName, "event", event)
				}
				buffer.WriteString(variableValue)
			}
			lastIndex = match[1]
		}
		buffer.WriteString(message[lastIndex:])
		return *models.NewMqttOutgoingMessage(topic, buffer.String())
	}
}

func newWSMessageBuilder(msgType, message string) func(triggers.Event) models.WSOutgoingMessage {
	return func(event triggers.Event) models.WSOutgoingMessage {
		var buffer bytes.Buffer
		re := regexp.MustCompile(`{{(.*?)}}`)

		prefix := "{\"type\":\"" + msgType + "\", \"payload\":"
		suffix := "}"

		lastIndex := 0
		for _, match := range re.FindAllStringSubmatchIndex(message, -1) {
			buffer.WriteString(message[lastIndex:match[0]])

			variableName := message[match[2]:match[3]]

			switch {
			case strings.HasPrefix(variableName, "$userVariable."):
				variableValue := getUserVariable(variableName)
				buffer.WriteString(variableValue)
			case strings.HasPrefix(variableName, "$event"):
				variableValue, err := getEventVariable(variableName, event)
				if err != nil {
					slog.Error("error writing value from event to message", "error", err, "variable name", variableName, "event", event)
				}
				buffer.WriteString(variableValue)
			}
			lastIndex = match[1]
		}
		buffer.WriteString(message[lastIndex:])
		return models.WSOutgoingMessage{Type: msgType, Message: prefix + buffer.String() + suffix}
	}
}

func newUpdateMessageBuilder() func(triggers.Event) models.WSOutgoingMessage {
	return func(event triggers.Event) models.WSOutgoingMessage {
		message := UpdateMessage{
			SensorId: event.SenderId,
			Value:    event.Value,
		}

		payload, err := json.Marshal(message)
		if err != nil {
			slog.Error("error creating message for reading update", "error", err)
			return models.WSOutgoingMessage{Type: "update", Message: ""}
		}

		prefix := "{\"type\":\"update\", \"payload\":"
		suffix := "}"

		return models.WSOutgoingMessage{Type: "update", Message: prefix + string(payload) + suffix}
	}
}

type UpdateMessage struct {
	SensorId string `json:"sensorId"`
	Value    any    `json:"value"`
}