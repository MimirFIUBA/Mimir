package mimir

type Condition interface {
	Evaluate() bool
	SetNewValue(newValue SensorValue)
}

type ReceiveValueCondition struct {
	hasReceivedValue bool
}

func (c *ReceiveValueCondition) SetNewValue(newValue SensorValue) {
	c.hasReceivedValue = true
}

func (c *ReceiveValueCondition) Evaluate() bool {
	return c.hasReceivedValue
}

type GenericCondition struct {
	ReferenceValue SensorValue
	TestValue      SensorValue
	Operator       string
}

func (c *GenericCondition) Evaluate() bool {
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

func compareInt(leftValue, rightValue int, operator string) bool {
	switch operator {
	case ">":
		return leftValue > rightValue
	case "<":
		return leftValue < rightValue
	case ">=":
		return leftValue >= rightValue
	case "<=":
		return leftValue <= rightValue
	case "==":
		return leftValue == rightValue
	case "!=":
		return leftValue != rightValue
	default:
		panic("Bad operator type")
	}
}

func compareFloat(leftValue, rightValue float64, operator string) bool {
	switch operator {
	case ">":
		return leftValue > rightValue
	case "<":
		return leftValue < rightValue
	case ">=":
		return leftValue >= rightValue
	case "<=":
		return leftValue <= rightValue
	case "==":
		return leftValue == rightValue
	case "!=":
		return leftValue != rightValue
	default:
		panic("Bad operator type")
	}
}

func compareString(leftValue, rightValue string, operator string) bool {
	switch operator {
	case "==":
		return leftValue == rightValue
	case "!=":
		return leftValue != rightValue
	default:
		panic("Bad operator type")
	}
}

func (c *GenericCondition) SetNewValue(newValue SensorValue) {
	c.TestValue = newValue
}
