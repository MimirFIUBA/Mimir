package triggers

type Condition interface {
	Evaluate(Event) bool
	SetEvent(Event)
	GetEventId() string
	SetEventId(string)
	GetSenderId() string
	SetSenderId(string)
}
