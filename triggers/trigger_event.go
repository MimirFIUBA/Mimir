package triggers

type EventTrigger struct {
	ID               string
	Name             string
	Condition        Condition
	stringCondition  string
	IsActive         bool
	Actions          []Action
	observedSubjects []Subject
	isScheduled      bool
}

func NewEventTrigger(name string) *EventTrigger {
	defaultCondition := TrueCondition{}
	return &EventTrigger{Name: name, Condition: &defaultCondition}
}

func (t *EventTrigger) Update(event Event) {
	if t.IsActive && t.Condition.Evaluate(event) {
		if !t.isScheduled || event.Type == SCHEDULER_ACTIVE {
			for _, action := range t.Actions {
				action.Execute(event)
			}
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

func (t *EventTrigger) AddAction(a Action, _ TriggerOptions) {
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

func (t *EventTrigger) UpdateActions(actions []Action, _ TriggerOptions) error {
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

func (t *EventTrigger) Activate() {
	t.IsActive = true
}

func (t *EventTrigger) Deactivate() {
	t.IsActive = false
}

func (t *EventTrigger) GetType() TriggerType {
	return EVENT_TRIGGER
}

func (t *EventTrigger) SetStatus(active bool) {
	if t.IsActive != active {
		if active {
			t.Activate()
		} else {
			t.Deactivate()
		}
	}
}

func (t *EventTrigger) SetScheduled(scheduled bool) {
	t.isScheduled = scheduled
}
