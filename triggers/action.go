package triggers

type Action interface {
	Execute(event Event)
}

type ActionEventType int

const (
	ACTIVE ActionEventType = iota
	INACTIVE
)
