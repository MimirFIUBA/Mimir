package triggers

import (
	"github.com/google/uuid"
)

type Trigger struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Condition Condition `json:"condition"`
	Actions   []Action  `json:"actions"`
}

func NewTrigger(name string) *Trigger {
	return &Trigger{uuid.New().String(), name, &TrueCondition{}, nil}
}

func (t *Trigger) Update(event Event) {
	t.Condition.SetEvent(event)
	if t.Condition.Evaluate() {
		for _, action := range t.Actions {
			action.Execute(event)
		}
	}
}

func (t *Trigger) GetID() string {
	return t.ID
}
