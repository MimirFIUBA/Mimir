package triggers

import (
	"os/exec"
)

type CommandAction struct {
	Command     string
	CommandArgs string
	NextAction  Action
}

func (a *CommandAction) Execute(event Event) {
	cmd := exec.Command(a.Command, a.CommandArgs)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return //TODO: see if we can do something with error
	}
	//TODO: see if we can do something with the output

	if a.NextAction != nil {
		a.NextAction.Execute(event)
	}
}
