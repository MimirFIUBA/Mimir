package triggers

import "time"

type Event struct {
	Name      string
	Timestamp time.Time
	Data      interface{}
	Id        string
	SenderId  string
}

func NewEvent() *Event {
	return &Event{}
}

// NewFloatEvent is just an empty event with only float data
func NewFloatEvent(data float64) *Event {
	return &Event{"", time.Now(), data, "", ""}
}

func (e *Event) GetId() string {
	return e.Id
}

func (e *Event) MatchesCondition(condition Condition) bool {
	return ((condition.GetEventId() != "" && e.Id == condition.GetEventId()) ||
		(condition.GetSenderId() != "" && e.SenderId == condition.GetSenderId()) ||
		(condition.GetEventId() == "" && condition.GetSenderId() == ""))
}
