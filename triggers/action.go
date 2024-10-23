package triggers

type Action interface {
	Execute(event Event)
}
