package triggers

import "fmt"

type SendMessageThroughChannel struct {
	MessageContructor       func(Event) string
	Message                 string
	OutgoingMessagesChannel chan string
}

func NewSendMessageThroughChannel(channel chan string) *SendMessageThroughChannel {
	return &SendMessageThroughChannel{}
}

func (action *SendMessageThroughChannel) Execute(event Event) {
	fmt.Println("Execute")
	if action.MessageContructor != nil {
		action.OutgoingMessagesChannel <- action.MessageContructor(event)
	} else {
		action.OutgoingMessagesChannel <- action.Message
	}
	fmt.Println("End Execute")
}
