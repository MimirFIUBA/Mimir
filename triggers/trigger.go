package triggers

import (
	"github.com/google/uuid"
)

type Trigger struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Condition Condition `json:"condition"`
	IsActive  bool      `json:"active"`
	Actions   []Action  `json:"actions"`
}

func NewTrigger(name string) *Trigger {
	defaultCondition := TrueCondition{}
	return &Trigger{ID: uuid.New().String(), Name: name, Condition: &defaultCondition}
}

func (t *Trigger) Update(event Event) {
	if t.Condition.Evaluate(event) {
		for _, action := range t.Actions {
			action.Execute(event)
		}
	}
}

func (t *Trigger) GetID() string {
	return t.ID
}

func (t *Trigger) SetCondition(c Condition) {
	t.Condition = c
}

func (t *Trigger) AddAction(a Action) {
	t.Actions = append(t.Actions, a)
}
