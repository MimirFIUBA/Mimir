package triggers

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

func (c *BetweenCondition) String() string {
	return "BetweenCondition"
}
