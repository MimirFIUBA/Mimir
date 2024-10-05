package triggers

import (
	"fmt"

	"github.com/google/uuid"
)

type Trigger struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Condition Condition `json:"condition"`
	Actions   []Action  `json:"actions"`
}

func NewTrigger(name string) *Trigger {
	defaultCondition := TrueCondition{}
	return &Trigger{uuid.New().String(), name, &defaultCondition, nil}
}

func (t *Trigger) Update(event Event) {
	if t.Condition.Evaluate(event) {
		fmt.Println("Actions: ", t.Actions)
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
