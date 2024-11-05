package triggers

import (
	"fmt"
	"log"
	"os/exec"
)

type CommandAction struct {
	Command     string
	CommandArgs string
	NextAction  Action
}

func (a *CommandAction) Execute(event Event) {
	cmd := exec.Command(a.Command, a.CommandArgs)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)

	if a.NextAction != nil {
		a.NextAction.Execute(event)
	}
}
