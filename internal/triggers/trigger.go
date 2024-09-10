package triggers

import (
	"github.com/google/uuid"
)

type TriggerEngine struct{}

type Trigger struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Condition Condition `json:"condition"`
	Actions   []Action  `json:"actions"`
}

func NewTrigger(name string) *Trigger {
	return &Trigger{uuid.New().String(), name, nil, nil}
}

func (t *Trigger) Update(event Event) {
	t.Condition.SetEvent(event)
	if t.Condition.Evaluate() {
		for _, action := range t.Actions {
			action.Execute()
		}
	}
}

func (t *Trigger) GetID() string {
	return t.ID
}
