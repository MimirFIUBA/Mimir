package triggers

import (
	"fmt"
	"math"
	"reflect"
)

type Condition interface {
	Evaluate() bool
	SetEvent(Event)
	GetEventId() string
	SetEventId(string)
}

type TrueCondition struct{}

func (c *TrueCondition) GetEventId() string {
	return ""
}

func (c *TrueCondition) SetEventId(id string) {}

func (c *TrueCondition) SetEvent(event Event) {}

func (c *TrueCondition) Evaluate() bool {
	return true
}

// Receive value condition: activates once it receives a value
type ReceiveValueCondition struct {
	valueId          string
	hasReceivedValue bool
}

func (c *ReceiveValueCondition) GetEventId() string {
	return c.valueId
}

func (c *ReceiveValueCondition) SetEventId(id string) {
	c.valueId = id
}

func (c *ReceiveValueCondition) SetEvent(event Event) {
	c.hasReceivedValue = true
}

func (c *ReceiveValueCondition) Evaluate() bool {
	return c.hasReceivedValue
}

// Compare condition, compare a value from an event to another reference value.
type CompareCondition struct {
	Operator       string
	ReferenceValue interface{}
	TestValue      interface{}
	valueId        string
}

func (c *CompareCondition) GetEventId() string {
	return c.valueId
}

func (c *CompareCondition) SetEventId(id string) {
	c.valueId = id
}

func (c *CompareCondition) Evaluate() bool {
	switch rightValue := c.ReferenceValue.(type) {
	case int:
		testValue := c.TestValue.(int)
		return compareInt(testValue, rightValue, c.Operator)
	case float64:
		testValue := c.TestValue.(float64)
		return compareFloat(testValue, rightValue, c.Operator)
	case string:
		testValue := c.TestValue.(string)
		return compareString(testValue, rightValue, c.Operator)
	default:
		panic("Bad values to compare")
	}
}

func (c *CompareCondition) SetEvent(event Event) {
	c.TestValue = event.Data
}

// And condition
type AndCondition struct {
	Conditions []Condition
}

func (c *AndCondition) GetEventId() string {
	return ""
}

func (c *AndCondition) SetEventId(id string) {}

func (c *AndCondition) Evaluate() bool {
	for _, condition := range c.Conditions {
		if !condition.Evaluate() {
			return false
		}
	}
	return true
}

func (c *AndCondition) SetEvent(event Event) {
	for _, condition := range c.Conditions {
		if condition.GetEventId() == event.GetId() || condition.GetEventId() == "" {
			condition.SetEvent(event)
		}
	}
}

// OR Condition
type OrCondition struct {
	Conditions []Condition
}

func (c *OrCondition) GetEventId() string {
	return ""
}

func (c *OrCondition) SetEventId(id string) {}

func (c *OrCondition) Evaluate() bool {
	for _, condition := range c.Conditions {
		if condition.Evaluate() {
			return true
		}
	}
	return false
}

func (c *OrCondition) SetEvent(event Event) {
	for _, condition := range c.Conditions {
		if condition.GetEventId() == event.GetId() || condition.GetEventId() == "" {
			condition.SetEvent(event)
		}
	}
}

// Custom condition
type CustomCondition struct {
	eventId      string
	currentEvent Event
	EvalFunc     func(event Event) bool
}

func (c *CustomCondition) Evaluate() bool {
	return c.EvalFunc(c.currentEvent)
}

func (c *CustomCondition) SetEvent(event Event) {
	c.currentEvent = event
}
func (c *CustomCondition) GetEventId() string {
	return c.eventId
}

func (c *CustomCondition) SetEventId(eventId string) {
	c.eventId = eventId
}

// Between condition
// Condición de rango (Between)
type BetweenCondition struct {
	CurrentValue interface{}
	Min          interface{}
	Max          interface{}
	eventId      string
}

func (c *BetweenCondition) Evaluate() bool {
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

func (c *BetweenCondition) SetEvent(e Event) {
	c.CurrentValue = e.Data
}

func (c *BetweenCondition) GetEventId() string {
	return c.eventId
}

func (c *BetweenCondition) SetEventId(id string) {
	c.eventId = id
}

// Condición de variación (Delta)
type DeltaCondition struct {
	CurrentValue  interface{}
	PreviousValue interface{}
	Threshold     interface{}
	eventId       string
}

func (c *DeltaCondition) Evaluate() bool {
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
	c.PreviousValue = c.CurrentValue
	c.CurrentValue = e.Data
}

func (c *DeltaCondition) GetEventId() string {
	return c.eventId
}

func (c *DeltaCondition) SetEventId(id string) {
	c.eventId = id
}

// Condición de promedio (Average)
// TODO: agregar time frame
type AverageCondition struct {
	Values       []interface{}
	Threshold    interface{}
	MinValueSize int
	MaxValueSize int
	eventId      string
	nextValueIdx int
	valueCount   int
}

func (c *AverageCondition) Evaluate() bool {
	if len(c.Values) < c.MinValueSize {
		return false
	}

	var sum float64
	for _, v := range c.Values {
		switch v := v.(type) {
		case int:
			sum += float64(v)
		case float64:
			sum += v
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(v)))
		}
	}

	avg := sum / float64(c.valueCount)
	thresholdFloat := c.Threshold.(float64)

	return avg > thresholdFloat
}

func (c *AverageCondition) SetEvent(e Event) {
	if c.nextValueIdx > c.MaxValueSize {
		c.nextValueIdx = 0
	}
	c.Values[c.nextValueIdx] = e.Data
	c.nextValueIdx++
	c.valueCount++
	if c.valueCount > c.MaxValueSize {
		c.valueCount = c.MaxValueSize
	}
}

func (c *AverageCondition) GetEventId() string {
	return c.eventId
}

func (c *AverageCondition) SetEventId(id string) {
	c.eventId = id
}
