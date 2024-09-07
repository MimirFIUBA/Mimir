package triggers

type TriggerObserver interface {
	Update(Event)
	GetID() string
}

type Subject interface {
	register(observer TriggerObserver)
	deregister(observer TriggerObserver)
	notifyAll()
}
