package triggers

import "fmt"

type PrintAction struct {
	Name              string
	messageContructor func(Event) string
	Message           string
}

func NewPrintAction() *PrintAction {
	return &PrintAction{}
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
}
