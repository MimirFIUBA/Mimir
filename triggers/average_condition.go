package triggers

import (
	"fmt"
	"reflect"
	"time"
)

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

func (c *AverageCondition) SetCondition(condition Condition) {
	c.Condition = condition
}

func (c *AverageCondition) String() string {
	return fmt.Sprintf("AVG(%s)[%d, %d, %d] %v", c.senderId, c.MinSize, c.MaxSize, c.timeFrame, c.Condition)
}
