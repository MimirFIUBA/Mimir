package triggers

import "fmt"

type Action interface {
	Execute(event Event)
}

type PrintAction struct {
	Message string
}

func (action *PrintAction) Execute(event Event) {
	fmt.Println(action.Message)
}

type SendMessageThroughChannel struct {
	MessageContructor       func(Event) string
	Message                 string
	OutgoingMessagesChannel chan string
}

func (action *SendMessageThroughChannel) Execute(event Event) {
	if action.MessageContructor != nil {
		action.OutgoingMessagesChannel <- action.MessageContructor(event)
	} else {
		action.OutgoingMessagesChannel <- action.Message
	}
	// action.OutgoingMessagesChannel <- action.Message
}
