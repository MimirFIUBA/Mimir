package triggers

import (
	"fmt"
	"time"
)

// Condición de promedio (Average)
type AggregateCondition struct {
	EventBuffer       []Event
	Condition         Condition
	MinSize           int
	MaxSize           int
	AggregateFunction func(*AggregateCondition) float64
	eventId           string
	senderId          string
	end               int
	start             int
	eventCount        int
	timeFrame         time.Duration
}

func NewAggregateCondition(minSize, maxSize int, timeFrame time.Duration) *AggregateCondition {
	return &AggregateCondition{
		EventBuffer: make([]Event, maxSize),
		MinSize:     minSize,
		MaxSize:     maxSize,
		timeFrame:   timeFrame,
	}
}

func (c *AggregateCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	c.cleanBuffer()
	if c.eventCount < c.MinSize {
		return false
	}

	result := c.AggregateFunction(c)

	return c.Condition.Evaluate(*NewAggregateFunctionEvent(result))
}

func (c *AggregateCondition) cleanBuffer() {
	currentTime := time.Now()
	currentEvent := c.EventBuffer[c.start]
	for currentTime.Sub(currentEvent.Timestamp) > c.timeFrame && c.eventCount > 0 {
		c.start = (c.start + 1) % c.MaxSize
		c.eventCount--
		currentEvent = c.EventBuffer[c.start]
	}
}

func (c *AggregateCondition) SetEvent(e Event) {
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

func (c *AggregateCondition) GetEventId() string {
	return c.eventId
}

func (c *AggregateCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *AggregateCondition) GetSenderId() string {
	return c.senderId
}

func (c *AggregateCondition) SetSenderId(id string) {
	c.senderId = id
}

func (c *AggregateCondition) SetCondition(condition Condition) {
	c.Condition = condition
}

func (c *AggregateCondition) String() string {
	return fmt.Sprintf("AGGR(%s)[%d, %d, %d] %v", c.senderId, c.MinSize, c.MaxSize, c.timeFrame, c.Condition)
}
