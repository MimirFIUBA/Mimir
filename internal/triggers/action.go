package triggers

import "fmt"

type Action interface {
	Execute()
}

type PrintAction struct {
	Message string
}

func (action *PrintAction) Execute() {
	fmt.Println(action.Message)
}

type SendMessageThroughChannel struct {
	Message                 string
	OutgoingMessagesChannel chan string
}

func (action *SendMessageThroughChannel) Execute() {
	action.OutgoingMessagesChannel <- action.Message
}
