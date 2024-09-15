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
