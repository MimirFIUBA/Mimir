package triggers

import (
	"time"

	"github.com/google/uuid"
)

type TimeTrigger struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Condition    Condition     `json:"condition"`
	Actions      []Action      `json:"actions"`
	Duration     time.Duration `json:"duration"`
	ticker       *time.Ticker
	resetChannel chan bool
}

func NewTimeTrigger(name string, duration time.Duration) *TimeTrigger {
	return &TimeTrigger{uuid.New().String(), name, nil, nil, duration, time.NewTicker(duration), make(chan bool)}
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
