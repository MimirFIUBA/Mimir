package triggers

import "fmt"

type TrueCondition struct{}

func (c *TrueCondition) GetEventId() string {
	return ""
}

func (c *TrueCondition) SetEventId(id string) {}

func (c *TrueCondition) GetSenderId() string {
	return ""
}

func (c *TrueCondition) SetSenderId(id string) {}

func (c *TrueCondition) SetEvent(event Event) {}

func (c *TrueCondition) Evaluate(event Event) bool {
	return true
}

func (c *TrueCondition) String() string {
	return "TrueCondition"
}

// Compare condition, compare a value from an event to another reference value.
type CompareCondition struct {
	Operator       string
	ReferenceValue interface{}
	TestValue      interface{}
	eventId        string
	senderId       string
	hasTestValue   bool
}

func NewCompareCondition(operator string, referenceValue interface{}) *CompareCondition {
	return &CompareCondition{
		Operator:       operator,
		ReferenceValue: referenceValue,
		hasTestValue:   false,
	}
}

func (c *CompareCondition) GetEventId() string {
	return c.eventId
}

func (c *CompareCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *CompareCondition) GetSenderId() string {
	return c.senderId
}

func (c *CompareCondition) SetSenderId(id string) {
	c.senderId = id
}

func (c *CompareCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	if c.hasTestValue {
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
	return false
}

func (c *CompareCondition) SetEvent(event Event) {
	if event.MatchesCondition(c) {
		c.hasTestValue = true
		c.TestValue = event.Value
	}
}

func (c *CompareCondition) String() string {
	if c.senderId != "" {
		return fmt.Sprintf("$(%s) %s %v", c.senderId, c.Operator, c.ReferenceValue)
	}
	return fmt.Sprintf("%s %v", c.Operator, c.ReferenceValue)
}

// And condition
type AndCondition struct {
	Conditions []Condition
}

func NewAndCondition(conditions []Condition) *AndCondition {
	return &AndCondition{conditions}
}

func (c *AndCondition) GetEventId() string {
	return ""
}

func (c *AndCondition) SetEventId(id string) {}

func (c *AndCondition) GetSenderId() string {
	return ""
}

func (c *AndCondition) SetSenderId(id string) {}

func (c *AndCondition) Evaluate(event Event) bool {
	for _, condition := range c.Conditions {
		if !condition.Evaluate(event) {
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

func (c *AndCondition) String() string {
	return "AndCondition"
}

// OR Condition
type OrCondition struct {
	Conditions []Condition
}

func NewOrCondition(conditions []Condition) *OrCondition {
	return &OrCondition{conditions}
}

func (c *OrCondition) GetEventId() string {
	return ""
}

func (c *OrCondition) SetEventId(id string) {}

func (c *OrCondition) GetSenderId() string {
	return ""
}

func (c *OrCondition) SetSenderId(id string) {}

func (c *OrCondition) Evaluate(event Event) bool {
	for _, condition := range c.Conditions {
		if condition.Evaluate(event) {
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

func (c *OrCondition) String() string {
	a := ""
	for _, subCondition := range c.Conditions {
		a += fmt.Sprintf("%v ", subCondition)
	}

	return fmt.Sprintf("OrCondition: %s", a)
}

type NotCondition struct {
	Cond Condition
}

func (c *NotCondition) GetEventId() string {
	return ""
}

func (c *NotCondition) SetEventId(id string) {}

func (c *NotCondition) GetSenderId() string {
	return ""
}

func (c *NotCondition) SetSenderId(id string) {}

func (c *NotCondition) Evaluate(event Event) bool {
	return !c.Cond.Evaluate(event)
}

func (c *NotCondition) SetEvent(event Event) {
	c.Cond.SetEvent(event)
}

func (c *NotCondition) String() string {
	return "NotCondition"
}
