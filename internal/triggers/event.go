package triggers

import "time"

type Event struct {
	Name      string
	Timestamp time.Time
	Data      interface{}
	Id        string
}

func (e *Event) GetId() string {
	return e.Id
}
