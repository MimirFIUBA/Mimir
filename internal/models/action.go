package models

import "mimir/triggers"

type ActionFactory struct {
	outgoingMessageChannel chan string
	wsMessageChannel       chan string
}

func NewActionFactory(mqttMsgChan, wsMsgChan chan string) *ActionFactory {
	return &ActionFactory{mqttMsgChan, wsMsgChan}
}

func (f *ActionFactory) NewSendMQTTMessageAction(message string) triggers.SendMessageThroughChannel {
	return triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: f.outgoingMessageChannel}
}

func (f *ActionFactory) NewSendWebSocketMessageAction(message string) triggers.SendMessageThroughChannel {
	return triggers.SendMessageThroughChannel{
		Message:                 message,
		OutgoingMessagesChannel: f.wsMessageChannel}
}
