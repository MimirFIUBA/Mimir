package triggers

type ExecuteFunctionAction struct {
	Func       func(event Event, params map[string]interface{}) Event
	Params     map[string]interface{}
	NextAction Action
}

func NewExecuteFunctionAction(function func(event Event, params map[string]interface{})) *ExecuteFunctionAction {
	return &ExecuteFunctionAction{}
}

func (a *ExecuteFunctionAction) Execute(event Event) {
	nextEvent := a.Func(event, a.Params)

	if a.NextAction != nil {
		a.NextAction.Execute(nextEvent)
	}
}
