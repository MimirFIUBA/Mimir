package mimir

import (
	"fmt"
	"time"
)

type Trigger struct {
	Condition Condition `json:"condition"`
	Actions   []Action  `json:"actions"`
}

func (t *Trigger) Execute(newValue SensorReading) {
	t.Condition.SetNewValue(newValue.Value)
	if t.Condition.Evaluate() {
		for _, action := range t.Actions {
			action.Execute()
		}
	}
}

type TimeTrigger struct {
	Condition    Condition
	Actions      []Action
	Duration     time.Duration
	timer        *time.Ticker
	resetChannel chan bool
}

func NewTimeTrigger(condition Condition, actions []Action, duration time.Duration) *TimeTrigger {
	return &TimeTrigger{condition, actions, duration, time.NewTicker(duration), make(chan bool)}
}

func (t *TimeTrigger) Start() {
	// t.timer = time.NewTicker(t.Duration)
	// go func() {
	// 	t.Execute()
	// }()
	go func() {
		for {
			select {
			case <-t.resetChannel:
				t.Reset()
			case <-t.timer.C:
				t.Execute()
			}
		}
	}()
}

func (t *TimeTrigger) Reset() {
	fmt.Printf("%v - [RESET]\n", time.Now())
	if t.timer != nil {
		t.timer.Reset(t.Duration)
	}
}

func (t *TimeTrigger) Execute() {
	for _, action := range t.Actions {
		fmt.Printf("%v - [ACTION EXECUTE]\n", time.Now())
		action.Execute()
	}
}

func (t *TimeTrigger) Evaluate(newValue SensorReading) {
	if t.Condition != nil {
		t.Condition.SetNewValue(newValue.Value)
		if t.Condition.Evaluate() {
			fmt.Println("Evaluate")
			t.resetChannel <- true
		}
	} else {
		t.resetChannel <- true
	}
}
