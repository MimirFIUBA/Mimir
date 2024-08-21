package trigger

type Observer interface {
	Update(Event)
	GetID() string
}

type Subject interface {
	register(observer Observer)
	deregister(observer Observer)
	notifyAll()
}
