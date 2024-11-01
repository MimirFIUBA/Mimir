package triggers

type TriggerType int

const (
	EVENT_TRIGGER TriggerType = iota
	SWITCH_TRIGGER
	TIMER_TRIGGER
	FREQUENCY_TRIGGER
)

type Trigger interface {
	Update(Event)
	GetID() string
	SetID(string)
	UpdateCondition(string) error
	UpdateActions([]Action, TriggerOptions) error
	AddAction(Action, TriggerOptions)
	AddSubject(Subject)
	StopWatching()
	Activate()
	Deactivate()
	GetType() TriggerType
	SetStatus(bool)
	SetScheduled(bool)
}

type Subject interface {
	Register(observer Trigger)
	Deregister(observer Trigger)
	NotifyAll()
}

type TriggerOptions struct {
	ActionsEventType ActionEventType
}
