package triggers

import (
	"math"
)

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

func (c *DeltaCondition) String() string {
	return "DeltaCondition"
}
