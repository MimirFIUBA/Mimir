package factories

import (
	"mimir/internal/db"
	"mimir/triggers"
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
		Func: func(event triggers.Event, params map[string]interface{}) {
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
		},
		Params: params,
	}
}

func (f *ActionFactory) NewCommandAction(command string, args string) *triggers.CommandAction {
	return &triggers.CommandAction{Command: command, CommandArgs: args}

}
