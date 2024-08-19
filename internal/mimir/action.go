package mimir

import "fmt"

type Action interface {
	Execute()
}

type PrintAction struct{}

func (action *PrintAction) Execute() {
	fmt.Println("Action executed")
}

type SendMQTTMessageAction struct {
	Message                 string
	OutgoingMessagesChannel chan string
}

func (action *SendMQTTMessageAction) Execute() {
	action.OutgoingMessagesChannel <- action.Message
}
