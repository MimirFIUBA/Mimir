package triggers

type SendMessageThroughChannel struct {
	MessageContructor       func(Event) string
	Message                 string
	OutgoingMessagesChannel chan string
}

func NewSendMessageThroughChannel(channel chan string) *SendMessageThroughChannel {
	return &SendMessageThroughChannel{}
}

func (action *SendMessageThroughChannel) Execute(event Event) {
	if action.MessageContructor != nil {
		action.OutgoingMessagesChannel <- action.MessageContructor(event)
	} else {
		action.OutgoingMessagesChannel <- action.Message
	}
}
