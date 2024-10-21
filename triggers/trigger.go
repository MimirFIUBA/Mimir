package triggers

type Trigger struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Condition        Condition `json:"condition"`
	stringCondition  string
	IsActive         bool     `json:"active"`
	Actions          []Action `json:"actions"`
	observedSubjects []Subject
}

func NewTrigger(name string) *Trigger {
	defaultCondition := TrueCondition{}
	return &Trigger{Name: name, Condition: &defaultCondition}
}

func (t *Trigger) Update(event Event) {
	if t.IsActive && t.Condition.Evaluate(event) {
		for _, action := range t.Actions {
			action.Execute(event)
		}
	}
}

func (t *Trigger) GetID() string {
	return t.ID
}

func (t *Trigger) SetID(id string) {
	t.ID = id
}

func (t *Trigger) SetCondition(c Condition) {
	t.Condition = c
}

func (t *Trigger) AddAction(a Action) {
	t.Actions = append(t.Actions, a)
}

func (t *Trigger) GetConditionAsString() string {
	return t.stringCondition
}

func (t *Trigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	t.stringCondition = newCondition
	return nil
}

func (t *Trigger) UpdateActions(actions []Action) error {
	t.Actions = actions
	return nil
}

func (t *Trigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *Trigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}
