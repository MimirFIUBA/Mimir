package triggers

import (
	"time"

	"github.com/google/uuid"
)

type TimeTrigger struct {
	ID               string        `json:"id"`
	Name             string        `json:"name"`
	IsActive         bool          `json:"active"`
	Condition        Condition     `json:"condition"`
	Actions          []Action      `json:"actions"`
	Duration         time.Duration `json:"duration"`
	ticker           *time.Ticker
	resetChannel     chan bool
	observedSubjects []Subject
}

func NewTimeTrigger(name string, duration time.Duration) *TimeTrigger {
	return &TimeTrigger{
		ID:           uuid.New().String(),
		Name:         name,
		Duration:     duration,
		ticker:       time.NewTicker(duration),
		resetChannel: make(chan bool),
	}
}

func (t *TimeTrigger) Start() {
	go func() {
		for {
			select {
			case <-t.resetChannel:
				t.reset()
			case <-t.ticker.C:
				t.execute()
			}
		}
	}()
}

func (t *TimeTrigger) reset() {
	if t.ticker != nil {
		t.ticker.Reset(t.Duration)
	}
}

func (t *TimeTrigger) execute() {
	for _, action := range t.Actions {
		action.Execute(*NewEvent())
	}
}

func (t *TimeTrigger) evaluate(event Event) {
	if t.Condition != nil {
		if t.Condition.Evaluate(event) {
			t.resetChannel <- true
		}
	} else {
		t.resetChannel <- true
	}
}

func (t *TimeTrigger) Update(event Event) {
	t.evaluate(event)
}

func (t *TimeTrigger) GetID() string {
	return t.ID
}

func (t *TimeTrigger) SetID(id string) {
	t.ID = id
}

func (t *TimeTrigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	return nil
}

func (t *TimeTrigger) UpdateActions(actions []Action) error {
	t.Actions = actions
	return nil
}

func (t *TimeTrigger) AddAction(a Action) {
	t.Actions = append(t.Actions, a)
}

func (t *TimeTrigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *TimeTrigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}
