package triggers

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

// Compare condition, compare a value from an event to another reference value.
type CompareCondition struct {
	Operator       string
	ReferenceValue interface{}
	TestValue      interface{}
	eventId        string
	senderId       string
}

func NewCompareCondition(operator string, referenceValue interface{}) *CompareCondition {
	//TODO: see what to do with first test value (set to reference value)
	return &CompareCondition{operator, referenceValue, referenceValue, "", ""}
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
	if event.MatchesCondition(c) {
		c.TestValue = event.Data
	}
}

// And condition
type AndCondition struct {
	Conditions []Condition
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

// OR Condition
type OrCondition struct {
	Conditions []Condition
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
