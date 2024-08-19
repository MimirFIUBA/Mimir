package mimir

type Trigger struct {
	Condition Condition `json:"condition"`
	Actions   []Action  `json:"action"`
}

func (t *Trigger) Execute(newValue SensorReading) {
	t.Condition.SetNewValue(newValue.Value)
	if t.Condition.Evaluate() {
		for _, action := range t.Actions {
			action.Execute()
		}
	}
}
