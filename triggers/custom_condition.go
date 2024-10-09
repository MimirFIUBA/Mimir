package triggers

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
