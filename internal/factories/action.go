package factories

import (
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/models"
	"mimir/triggers"
	"reflect"
	"strings"
	"time"
)

type ActionFactory struct {
	outgoingMessageChannel chan models.MqttOutgoingMessage
	wsMessageChannel       chan models.WSOutgoingMessage
}

type ActionType int

const (
	PRINT_ACTION ActionType = iota
	MQTT_ACTION
	WS_ACTION
)

const (
	USER_VARIABLE_PREFIX = "$userVariable"
	EVENT_PREFIX         = "$event"
)

func NewActionFactory(mqttMsgChan chan models.MqttOutgoingMessage, wsMsgChan chan models.WSOutgoingMessage) *ActionFactory {
	return &ActionFactory{mqttMsgChan, wsMsgChan}
}

func (f *ActionFactory) NewSendMQTTMessageAction(topic, message string) *triggers.SendMessageThroughChannel[models.MqttOutgoingMessage] {
	return &triggers.SendMessageThroughChannel[models.MqttOutgoingMessage]{
		Message:                 *models.NewMqttOutgoingMessage(topic, message),
		MessageContructor:       newMqttMessageBuilder(topic, message),
		OutgoingMessagesChannel: f.outgoingMessageChannel,
	}
}

func (f *ActionFactory) NewSendWebSocketMessageAction(message string) *triggers.SendMessageThroughChannel[models.WSOutgoingMessage] {
	return &triggers.SendMessageThroughChannel[models.WSOutgoingMessage]{
		Message:                 models.WSOutgoingMessage{Message: message},
		OutgoingMessagesChannel: f.wsMessageChannel,
	}
}

func (f *ActionFactory) NewWebSocketUpdateMessageAction(message string) *triggers.SendMessageThroughChannel[models.WSOutgoingMessage] {
	return &triggers.SendMessageThroughChannel[models.WSOutgoingMessage]{
		Message:                 models.WSOutgoingMessage{Type: "update", Message: message},
		MessageContructor:       newUpdateMessageBuilder(),
		OutgoingMessagesChannel: f.wsMessageChannel,
	}
}

func (f *ActionFactory) NewChangeTriggerStatus(triggerName string, status bool) *triggers.ExecuteFunctionAction {
	params := map[string]interface{}{"name": triggerName, "status": status}
	return &triggers.ExecuteFunctionAction{
		Func: func(event triggers.Event, params map[string]interface{}) triggers.Event {
			triggerNameValue, exists := params["name"]
			if exists {
				triggerName, ok := triggerNameValue.(string)
				if ok {
					trigger := db.GetTriggerByName(triggerName)
					if trigger != nil {
						statusValue, exists := params["status"]
						if exists {
							status, ok := statusValue.(bool)
							if ok {
								trigger.SetStatus(status)
							}
						}
					}
				}
			}
			return event
		},
		Params: params,
	}
}

func (f *ActionFactory) NewCommandAction(command string, args string) *triggers.CommandAction {
	return &triggers.CommandAction{Command: command, CommandArgs: args}
}

func (f *ActionFactory) NewAlertMessageAction(message string) *triggers.ExecuteFunctionAction {
	params := map[string]interface{}{"message": message}
	actionWS := f.NewSendWebSocketMessageAction(message)
	actionWS.MessageContructor = newWSMessageBuilder("alert", message)
	actionMqtt := f.NewSendMQTTMessageAction(consts.MQTT_ALERT_TOPIC, message)
	actionMqtt.NextAction = actionWS
	actionCreateMessage := &triggers.ExecuteFunctionAction{
		Func: func(event triggers.Event, params map[string]interface{}) triggers.Event {

			additionalDetails := make(map[string]interface{})
			additionalDetails["senderId"] = event.SenderId
			additionalDetails["value"] = event.Value
			additionalDetails["event"] = event

			messageValue, exists := params["message"]
			if exists {
				message, ok := messageValue.(string)
				if ok {
					message := &models.Message{
						Type:              models.ALERT_MESSAGE_TYPE,
						Body:              message,
						AdditionalDetails: additionalDetails,
						CreatedDate:       time.Now(),
					}
					_, err := db.Database.InsertMessage(message)
					if err != nil {
						slog.Error("error inserting message", "error", err)
					}
				}
			}
			return event
		},
		Params:     params,
		NextAction: actionMqtt,
	}
	return actionCreateMessage
}

func getUserVariable(variableName string) string {
	variableName = variableName[14:]
	userVariable, exists := db.GetUserVariable(variableName)
	if exists {
		return fmt.Sprintf("%v", userVariable.Value)
	} else {
		return fmt.Sprintf("{{%s}}", variableName)
	}
}

func getEventVariable(variableName string, event triggers.Event) (string, error) {
	parts := strings.Split(variableName[len("$event."):], ".")

	// Usa reflección para acceder a los campos del evento
	var currentValue interface{} = event
	for _, part := range parts {
		if part == "reading" {
			dataMap, ok := event.Data.(map[string]interface{})
			if !ok {
				return "", fmt.Errorf("event data is not a map")
			}
			reading, ok := dataMap["reading"]
			if !ok {
				return "", fmt.Errorf("event data has no reading")
			}
			currentValue = reading
			continue
		}

		val := reflect.ValueOf(currentValue)

		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() == reflect.Struct {
			field := val.FieldByName(part)
			if !field.IsValid() {
				return "", fmt.Errorf("field %s not found in event", part)
			}
			currentValue = field.Interface()

		} else if val.Kind() == reflect.Map {
			key := reflect.ValueOf(part)
			field := val.MapIndex(key)
			if !field.IsValid() {
				return "", fmt.Errorf("key %s not found in map", part)
			}
			currentValue = field.Interface()
		} else {
			return "", fmt.Errorf("invalid type encountered in path: %s", part)
		}
	}

	return fmt.Sprintf("%v", currentValue), nil
}
