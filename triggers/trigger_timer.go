package triggers

import (
	"time"

	"github.com/google/uuid"
)

type TimerTrigger struct {
	ID               string
	Name             string
	IsActive         bool
	Condition        Condition
	Actions          []Action
	Timeout          time.Duration
	ticker           *time.Ticker
	resetChannel     chan bool
	observedSubjects []Subject
}

func NewTimerTrigger(name string, timeout time.Duration) *TimerTrigger {
	return &TimerTrigger{
		ID:           uuid.New().String(),
		Name:         name,
		Timeout:      timeout,
		ticker:       time.NewTicker(timeout),
		resetChannel: make(chan bool),
	}
}

func (t *TimerTrigger) Start() {
	go func() {
		//TODO: ADD CANCEL TO THIS
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

func (t *TimerTrigger) reset() {
	if t.ticker != nil {
		t.ticker.Reset(t.Timeout)
	}
}

func (t *TimerTrigger) execute() {
	if t.IsActive {
		for _, action := range t.Actions {
			action.Execute(*NilEvent())
		}
	}
}

func (t *TimerTrigger) evaluate(event Event) {
	if t.Condition != nil {
		if t.Condition.Evaluate(event) {
			t.resetChannel <- true
		}
	} else {
		t.resetChannel <- true
	}
}

func (t *TimerTrigger) Update(event Event) {
	t.evaluate(event)
}

func (t *TimerTrigger) GetID() string {
	return t.ID
}

func (t *TimerTrigger) SetID(id string) {
	t.ID = id
}

func (t *TimerTrigger) UpdateCondition(newCondition string) error {
	condition, err := BuildConditionFromString(newCondition)
	if err != nil {
		return err
	}
	t.Condition = condition
	return nil
}

func (t *TimerTrigger) UpdateActions(actions []Action) error {
	t.Actions = actions
	return nil
}

func (t *TimerTrigger) AddAction(a Action) {
	t.Actions = append(t.Actions, a)
}

func (t *TimerTrigger) AddSubject(subject Subject) {
	t.observedSubjects = append(t.observedSubjects, subject)
}

func (t *TimerTrigger) StopWatching() {
	for _, subject := range t.observedSubjects {
		subject.Deregister(t)
	}
}

func (t *TimerTrigger) Activate() {
	t.IsActive = true
	t.Start()
}

func (t *TimerTrigger) Deactivate() {
	t.IsActive = false
}

func (t *TimerTrigger) GetType() TriggerType {
	return TIMER_TRIGGER
}

func (t *TimerTrigger) SetStatus(active bool) {
	if t.IsActive != active {
		if active {
			t.Activate()
		} else {
			t.Deactivate()
		}
	}
}

func (t *TimerTrigger) UpdateTimeout(newTimeout time.Duration) {
	t.Timeout = newTimeout
	t.reset()
}
