package triggers

type TriggerObserver interface {
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
	Register(observer TriggerObserver)
	Deregister(observer TriggerObserver)
	NotifyAll()
}
