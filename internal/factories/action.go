package factories

import (
	"bytes"
	"fmt"
	"log/slog"
	"mimir/internal/consts"
	"mimir/internal/db"
	"mimir/internal/models"
	"mimir/triggers"
	"regexp"
	"time"
)

type ActionFactory struct {
	outgoingMessageChannel chan models.MqttOutgoingMessage
	wsMessageChannel       chan string
}

type ActionType int

const (
	PRINT_ACTION ActionType = iota
	MQTT_ACTION
	WS_ACTION
)

func NewActionFactory(mqttMsgChan chan models.MqttOutgoingMessage, wsMsgChan chan string) *ActionFactory {
	return &ActionFactory{mqttMsgChan, wsMsgChan}
}

func (f *ActionFactory) NewSendMQTTMessageAction(topic, message string) *triggers.SendMessageThroughChannel[models.MqttOutgoingMessage] {
	var msgConstructor = func(event triggers.Event) models.MqttOutgoingMessage {
		var buffer bytes.Buffer
		re := regexp.MustCompile(`{{(.*?)}}`)

		lastIndex := 0
		for _, match := range re.FindAllStringSubmatchIndex(message, -1) {
			// Agregar la parte del mensaje sin reemplazar
			buffer.WriteString(message[lastIndex:match[0]])

			// Obtener el nombre de la variable
			variableName := message[match[2]:match[3]]
			// Reemplazar con el valor de la variable si existe
			userVariable, exists := db.GetUserVariable(variableName)
			if exists {
				userVariableStringer, ok := userVariable.Value.(fmt.Stringer)
				if ok {
					buffer.WriteString(userVariableStringer.String())
				}
			} else {
				buffer.WriteString("{{" + variableName + "}}") // Deja el placeholder si no existe
			}

			lastIndex = match[1]
		}
		buffer.WriteString(message[lastIndex:]) // Añadir el resto del mensaje
		return *models.NewMqttOutgoingMessage(topic, buffer.String())
	}

	return &triggers.SendMessageThroughChannel[models.MqttOutgoingMessage]{
		Message:                 *models.NewMqttOutgoingMessage(topic, message),
		MessageContructor:       msgConstructor,
		OutgoingMessagesChannel: f.outgoingMessageChannel}
}

func (f *ActionFactory) NewSendWebSocketMessageAction(message string) *triggers.SendMessageThroughChannel[string] {
	return &triggers.SendMessageThroughChannel[string]{
		Message:                 message,
		OutgoingMessagesChannel: f.wsMessageChannel}
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
