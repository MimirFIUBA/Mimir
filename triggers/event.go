package triggers

import (
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	SCHEDULER_ACTIVE EventType = iota
	NEW_READING
	STATISTIC_CALCULATION
	CHANNEL_MESSAGE_SENT
)

type Event struct {
	Name      string
	Timestamp time.Time
	Data      interface{}
	Value     interface{}
	Id        string
	SenderId  string
	Type      EventType
}

func NewEvent() *Event {
	return &Event{Id: uuid.New().String()}
}

func NilEvent() *Event {
	return &Event{Id: uuid.Nil.String()}
}

// NewFloatEvent is just an empty event with only float data
func NewFloatEvent(data float64) *Event {
	return &Event{
		Id:        uuid.New().String(),
		Timestamp: time.Now(),
		Type:      STATISTIC_CALCULATION,
		Data:      data,
	}
}

func (e *Event) GetId() string {
	return e.Id
}

func (e *Event) MatchesCondition(condition Condition) bool {
	return ((condition.GetEventId() != "" && e.Id == condition.GetEventId()) ||
		(condition.GetSenderId() != "" && e.SenderId == condition.GetSenderId()) ||
		(condition.GetEventId() == "" && condition.GetSenderId() == ""))
}
