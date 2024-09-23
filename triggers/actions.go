package triggers

import (
	"fmt"
	"log"
	"os/exec"
)

type Action interface {
	Execute(event Event)
}

type PrintAction struct {
	messageContructor func(Event) string
	Message           string
}

func NewPrintAction() *PrintAction {
	return &PrintAction{nil, ""}
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
}

type CommandAction struct {
	Command     string
	CommandArgs string
}

func (a *CommandAction) Execute(event Event) {
	cmd := exec.Command(a.Command, a.CommandArgs)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)
}
