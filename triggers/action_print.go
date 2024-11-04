package triggers

import "fmt"

type PrintAction struct {
	Name              string
	messageContructor func(Event) string
	Message           string
	NextAction        Action
}

func NewPrintAction(message string) *PrintAction {
	return &PrintAction{Message: message}
}

func (a *PrintAction) SetMessage(message string) {
	a.Message = message
}

func (a *PrintAction) SetMessageConstructor(messageConstructor func(Event) string) {
	a.messageContructor = messageConstructor
}

func (action *PrintAction) Execute(event Event) {
	message := action.Message
	if action.messageContructor != nil {
		message = action.messageContructor(event)
	}
	fmt.Println(message)
	if action.NextAction != nil {
		action.NextAction.Execute(event)
	}
}
