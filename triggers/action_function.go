package triggers

type ExecuteFunctionAction struct {
	Func   func(event Event, params map[string]interface{})
	Params map[string]interface{}
}

func (a *ExecuteFunctionAction) Execute(event Event) {
	a.Func(event, a.Params)
}
