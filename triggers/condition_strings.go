package triggers

import (
	"regexp"
)

type MatchesCondition struct {
	Pattern   *regexp.Regexp
	eventId   string
	testValue string
	senderId  string
}

func NewMatchesCondition(patternString string) *MatchesCondition {
	pattern := regexp.MustCompile(patternString)
	return &MatchesCondition{pattern, "", "", ""}
}

func (c *MatchesCondition) GetEventId() string {
	return c.eventId
}

func (c *MatchesCondition) SetEventId(id string) {
	c.eventId = id
}

func (c *MatchesCondition) GetSenderId() string {
	return c.senderId
}

func (c *MatchesCondition) SetSenderId(id string) {
	c.senderId = id
}

func (c *MatchesCondition) SetEvent(event Event) {
	if event.MatchesCondition(c) {
		switch data := event.Value.(type) {
		case string:
			c.testValue = data
		case []byte:
			c.testValue = string(data)
		default:
			panic("Wrong value type for matches condition")
		}
	}
}

func (c *MatchesCondition) Evaluate(event Event) bool {
	c.SetEvent(event)
	matched := c.Pattern.MatchString(c.testValue)
	return matched
}

func (c *MatchesCondition) String() string {
	return "MatchesCondition"
}
