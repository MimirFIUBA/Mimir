package mimir

type Condition interface {
	Evaluate() bool
	SetNewValue(newValue SensorValue)
}

type MaxValueCondition struct {
	MaxValue  SensorValue
	TestValue SensorValue
}

func (c *MaxValueCondition) Evaluate() bool {
	switch maxValue := c.MaxValue.(type) {
	case int:
		testValue := c.TestValue.(int)
		return testValue > maxValue
	case float64:
		testValue := c.TestValue.(float64)
		return testValue > maxValue
	default:
		panic("Bad values to compare")
	}
}

func (c *MaxValueCondition) SetNewValue(newValue SensorValue) {
	c.TestValue = newValue
}
