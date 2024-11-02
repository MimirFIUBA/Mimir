package triggers

type SwitchTrigger struct {
	ID               string
	Name             string
	Condition        Condition
	stringCondition  string
	IsActive         bool
	ActiveActions    []Action
	InactiveActions  []Action
	observedSubjects []Subject
	lastStatus       bool
	isScheduled      bool
}

func NewSwitchTrigger(name string) *SwitchTrigger {
	defaultCondition := TrueCondition{}
	return &SwitchTrigger{Name: name, Condition: &defaultCondition, lastStatus: false}
}

func (t *SwitchTrigger) Update(event Event) {
	if t.IsActive {
		condition := t.Condition.Evaluate(event)
		if t.lastStatus != condition {
			t.lastStatus = condition
			if condition {
				for _, action := range t.ActiveActions {
					action.Execute(event)
				}
			} else {
				for _, action := range t.InactiveActions {
					action.Execute(event)
				}

			}
		}
	}
}

func (t *SwitchTrigger) GetID() string {
	return t.ID
}

func (t *SwitchTrigger) SetID(id string) {
	t.ID = id
}

func (t *SwitchTrigger) SetCondition(c Condition) {
	t.Condition = c
}

func (t *SwitchTrigger) AddAction(a Action, opts TriggerOptions) {
	if opts.ActionsEventType == ACTIVE {
		t.ActiveActions = append(t.ActiveActions, a)
	} else if opts.ActionsEventType == INACTIVE {
		t.InactiveActions = append(t.InactiveActions, a)
	}
}

func (t *SwitchTrigger) UpdateActions(actions []Action, opts TriggerOptions) error {
	if opts.ActionsEventType == ACTIVE {
		t.ActiveActions = actions
	} else if opts.ActionsEventType == INACTIVE {
		t.InactiveActions = actions
	}
	return nil
}

func (t *SwitchTrigger) GetConditionAsString() string {
	return t.stringCondition
}

func (t *SwitchTrigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	t.stringCondition = newCondition
	return nil
}

func (t *SwitchTrigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *SwitchTrigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}

func (t *SwitchTrigger) Activate() {
	t.IsActive = true
}

func (t *SwitchTrigger) Deactivate() {
	t.IsActive = false
}

func (t *SwitchTrigger) GetType() TriggerType {
	return SWITCH_TRIGGER
}

func (t *SwitchTrigger) SetStatus(active bool) {
	if t.IsActive != active {
		if active {
			t.Activate()
		} else {
			t.Deactivate()
		}
	}
}

func (t *SwitchTrigger) SetScheduled(scheduled bool) {
	t.isScheduled = scheduled
}

func (t *SwitchTrigger) GetName() string {
	return t.Name
}
