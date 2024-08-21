package trigger

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

type SendMQTTMessageAction struct {
	Message                 string
	OutgoingMessagesChannel chan string
}

func (action *SendMQTTMessageAction) Execute() {
	action.OutgoingMessagesChannel <- action.Message
}
