package triggers

type Trigger interface {
	Update(Event)
	GetID() string
	SetID(string)
	UpdateCondition(string) error
	UpdateActions([]Action) error
	AddAction(Action)
	AddSubject(Subject)
	StopWatching()
}

type Subject interface {
	Register(observer Trigger)
	Deregister(observer Trigger)
	NotifyAll()
}

type TriggerType int

const (
	EVENT_TRIGGER TriggerType = iota
	TIMER_TRIGGER
	FREQUENCY_TRIGGER
)
