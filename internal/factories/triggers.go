package factories

import (
	"fmt"
	"mimir/internal/consts"
	"mimir/triggers"
	"time"
)

type TriggerFactory struct {
	actionFactory *ActionFactory
}

func NewTriggerFactory(actionFactory *ActionFactory) *TriggerFactory {
	return &TriggerFactory{actionFactory}
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
	case triggers.SWITCH_TRIGGER:
		return triggers.NewSwitchTrigger(opts.Name), nil
	case triggers.TIMER_TRIGGER:
		return triggers.NewTimerTrigger(opts.Name, opts.Timeout), nil
	case triggers.FREQUENCY_TRIGGER:
		return triggers.NewFrequencyTrigger(opts.Name, opts.Frequency), nil
	default:
		return nil, fmt.Errorf("wrong trigger type")
	}
}

func (f *TriggerFactory) BuildNewReadingNotificationTrigger() (triggers.Trigger, error) {
	trigger, err := f.BuildTrigger(TriggerOptions{
		Name:        "update",
		Frequency:   consts.TRIGGER_UPDATE_NOTIFICATION_FREQUENCY,
		TriggerType: triggers.FREQUENCY_TRIGGER,
	})
	if err != nil {
		return nil, err
	}

	action := f.actionFactory.NewWebSocketUpdateMessageAction("{\"type\":\"update\"}")
	trigger.AddAction(action, triggers.TriggerOptions{})

	return trigger, nil
}
