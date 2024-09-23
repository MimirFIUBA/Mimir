package triggers

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

type Condition interface {
	Evaluate(Event) bool
	SetEvent(Event)
	GetEventId() string
	SetEventId(string)
	GetSenderId() string
	SetSenderId(string)
}

// Receive value condition: activates once it receives a value
type ReceiveValueCondition struct {
	eventId          string
	senderId         string
	hasReceivedValue bool
}

func (c *ReceiveValueCondition) GetEventId() string {
	return c.eventId
}

func (c *ReceiveValueCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *ReceiveValueCondition) GetSenderId() string {
	return c.senderId
}

func (c *ReceiveValueCondition) SetSenderId(id string) {
	c.senderId = id
}

func (c *ReceiveValueCondition) SetEvent(event Event) {
	if event.MatchesCondition(c) {
		c.hasReceivedValue = true
	}
}

func (c *ReceiveValueCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	return c.hasReceivedValue
}

// Custom condition
type CustomCondition struct {
	currentEvent Event
	EvalFunc     func(event Event) bool
	eventId      string
	senderId     string
}

func (c *CustomCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	return c.EvalFunc(c.currentEvent)
}

func (c *CustomCondition) SetEvent(event Event) {
	if event.MatchesCondition(c) {
		c.currentEvent = event
	}
}

func (c *CustomCondition) GetEventId() string {
	return c.eventId
}

func (c *CustomCondition) SetEventId(eventId string) {
	c.eventId = eventId
}

func (c *CustomCondition) GetSenderId() string {
	return c.senderId
}

func (c *CustomCondition) SetSenderId(senderId string) {
	c.senderId = senderId
}

// Between condition
// Condición de rango (Between)
type BetweenCondition struct {
	CurrentValue interface{}
	Min          interface{}
	Max          interface{}
	eventId      string
	senderId     string
}

func (c *BetweenCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	switch testValue := c.CurrentValue.(type) {
	case int:
		min := c.Min.(int)
		max := c.Max.(int)
		return testValue >= min && testValue <= max
	case float64:
		min := c.Min.(float64)
		max := c.Max.(float64)
		return testValue >= min && testValue <= max
	default:
		panic("Bad values to compare")
	}
}

func (c *BetweenCondition) SetEvent(event Event) {
	if event.MatchesCondition(c) {
		c.CurrentValue = event.Data
	}
}

func (c *BetweenCondition) GetEventId() string {
	return c.eventId
}

func (c *BetweenCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *BetweenCondition) GetSenderId() string {
	return c.senderId
}

func (c *BetweenCondition) SetSenderId(id string) {
	c.senderId = id
}

// Condición de variación (Delta)
type DeltaCondition struct {
	CurrentValue  interface{}
	PreviousValue interface{}
	Threshold     interface{}
	eventId       string
	senderId      string
}

func (c *DeltaCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	switch currentValue := c.CurrentValue.(type) {
	case int:
		previousValue := c.PreviousValue.(int)
		threshold := c.Threshold.(int)
		delta := abs(currentValue - previousValue)
		return delta > threshold
	case float64:
		previousValue := c.PreviousValue.(float64)
		threshold := c.Threshold.(float64)
		delta := math.Abs(currentValue - previousValue)
		return delta > threshold
	default:
		panic("Bad values to compare")
	}
}

func (c *DeltaCondition) SetEvent(e Event) {
	if e.MatchesCondition(c) {
		c.PreviousValue = c.CurrentValue
		c.CurrentValue = e.Data
	}
}

func (c *DeltaCondition) GetEventId() string {
	return c.eventId
}

func (c *DeltaCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *DeltaCondition) GetSenderId() string {
	return c.senderId
}

func (c *DeltaCondition) SetSenderId(id string) {
	c.senderId = id
}

// Condición de promedio (Average)
type AverageCondition struct {
	EventBuffer []Event
	Condition   Condition
	MinSize     int
	MaxSize     int
	eventId     string
	senderId    string
	end         int
	start       int
	eventCount  int
	timeFrame   time.Duration
}

func NewAverageCondition(minSize, maxSize int, timeFrame time.Duration) *AverageCondition {
	return &AverageCondition{make([]Event, maxSize), nil, minSize, maxSize, "", "", 0, 0, 0, timeFrame}
}

func (c *AverageCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	c.cleanBuffer()
	if c.eventCount < c.MinSize {
		return false
	}

	avg := c.calculateAverage()
	fmt.Println("average", avg)

	return c.Condition.Evaluate(*NewFloatEvent(avg))
}

func (c *AverageCondition) calculateAverage() float64 {
	var sum float64
	for i := range c.eventCount {
		event := c.EventBuffer[(c.start+i)%c.MaxSize]
		data := event.Data
		switch data := data.(type) {
		case int:
			sum += float64(data)
		case float64:
			sum += data
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
		}
	}

	avg := sum / float64(c.eventCount)
	return avg
}

func (c *AverageCondition) cleanBuffer() {
	currentTime := time.Now()
	currentEvent := c.EventBuffer[c.start]
	for currentTime.Sub(currentEvent.Timestamp) > c.timeFrame && c.eventCount > 0 {
		c.start = (c.start + 1) % c.MaxSize
		c.eventCount--
		currentEvent = c.EventBuffer[c.start]
	}

}

func (c *AverageCondition) SetEvent(e Event) {
	if e.MatchesCondition(c) {
		c.EventBuffer[c.end] = e
		c.end = (c.end + 1) % c.MaxSize

		if c.eventCount == c.MaxSize {
			c.start = (c.start + 1) % c.MaxSize
		} else {
			c.eventCount++
		}
	}
}

func (c *AverageCondition) GetEventId() string {
	return c.eventId
}

func (c *AverageCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *AverageCondition) GetSenderId() string {
	return c.senderId
}

func (c *AverageCondition) SetSenderId(id string) {
	c.senderId = id
}
