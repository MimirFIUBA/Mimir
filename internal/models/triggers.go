package models

import (
	"fmt"
	"mimir/triggers"
	"time"
)

type TriggerFactory struct {
}

func NewTriggerFactory() *TriggerFactory {
	return &TriggerFactory{}
}

type TriggerOptions struct {
	Name        string
	Frequency   time.Duration
	Timeout     time.Duration
	TriggerType triggers.TriggerType
}

func (f *TriggerFactory) BuildTrigger(opts TriggerOptions) (triggers.Trigger, error) {
	switch opts.TriggerType {
	case triggers.EVENT_TRIGGER:
		return triggers.NewEventTrigger(opts.Name), nil
	case triggers.TIMER_TRIGGER:
		return triggers.NewTimerTrigger(opts.Name, opts.Timeout), nil
	case triggers.FREQUENCY_TRIGGER:
		return triggers.NewFrequencyTrigger(opts.Name, opts.Frequency), nil
	default:
		return nil, fmt.Errorf("wrong trigger type")
	}
}
