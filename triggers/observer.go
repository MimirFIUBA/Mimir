package triggers

type TriggerObserver interface {
	Update(Event)
	GetID() string
	SetID(string)
}

type Subject interface {
	Register(observer TriggerObserver)
	Deregister(observer TriggerObserver)
	NotifyAll()
}
