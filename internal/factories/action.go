package factories

import (
	"mimir/internal/db"
	"mimir/internal/models"
	"mimir/triggers"
	"time"
)

type ActionFactory struct {
	outgoingMessageChannel chan string
	wsMessageChannel       chan string
}

type ActionType int

const (
	PRINT_ACTION ActionType = iota
	MQTT_ACTION
	WS_ACTION
)

func NewActionFactory(mqttMsgChan, wsMsgChan chan string) *ActionFactory {
	return &ActionFactory{mqttMsgChan, wsMsgChan}
}

func (f *ActionFactory) NewSendMQTTMessageAction(message string) *triggers.SendMessageThroughChannel {
	return &triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: f.outgoingMessageChannel}
}

func (f *ActionFactory) NewSendWebSocketMessageAction(message string) *triggers.SendMessageThroughChannel {
	return &triggers.SendMessageThroughChannel{
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
	actionMqtt := f.NewSendMQTTMessageAction(message)
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
					db.Database.InsertMessage(message)
				}
			}
			return event
		},
		Params:     params,
		NextAction: actionMqtt,
	}
	return actionCreateMessage
}
