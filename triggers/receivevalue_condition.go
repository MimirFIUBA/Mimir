package triggers

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
