package triggers

type TriggerObserver interface {
	Update(Event)
	GetID() string
	SetID(string)
	UpdateCondition(string) error
}

type Subject interface {
	Register(observer TriggerObserver)
	Deregister(observer TriggerObserver)
	NotifyAll()
}
