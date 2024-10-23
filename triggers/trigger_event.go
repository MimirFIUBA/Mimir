package triggers

type EventTrigger struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Condition        Condition `json:"condition"`
	stringCondition  string
	IsActive         bool     `json:"active"`
	Actions          []Action `json:"actions"`
	observedSubjects []Subject
}

func NewEventTrigger(name string) *EventTrigger {
	defaultCondition := TrueCondition{}
	return &EventTrigger{Name: name, Condition: &defaultCondition}
}

func (t *EventTrigger) Update(event Event) {
	if t.IsActive && t.Condition.Evaluate(event) {
		for _, action := range t.Actions {
			action.Execute(event)
		}
	}
}

func (t *EventTrigger) GetID() string {
	return t.ID
}

func (t *EventTrigger) SetID(id string) {
	t.ID = id
}

func (t *EventTrigger) SetCondition(c Condition) {
	t.Condition = c
}

func (t *EventTrigger) AddAction(a Action) {
	t.Actions = append(t.Actions, a)
}

func (t *EventTrigger) GetConditionAsString() string {
	return t.stringCondition
}

func (t *EventTrigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	t.stringCondition = newCondition
	return nil
}

func (t *EventTrigger) UpdateActions(actions []Action) error {
	t.Actions = actions
	return nil
}

func (t *EventTrigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *EventTrigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}
