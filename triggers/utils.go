package triggers

import (
	"fmt"
	"reflect"
)

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

func abs(value int) int {
	if value >= 0 {
		return value
	}
	return -value
}

func calculateAverage(c *AggregateCondition) float64 {
	sum := calculateSum(c)
	avg := sum / float64(c.eventCount)
	return avg
}

func calculateCount(c *AggregateCondition) float64 {
	return float64(c.eventCount)
}

func calculateMin(c *AggregateCondition) float64 {
	var min float64
	firstEvent := c.EventBuffer[c.start]
	switch data := firstEvent.Value.(type) {
	case int:
		min = float64(data)
	case float64:
		min = data
	default:
		panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
	}

	for i := range c.eventCount {
		event := c.EventBuffer[(c.start+i)%c.MaxSize]
		data := event.Value
		switch data := data.(type) {
		case int:
			if min > float64(data) {
				min = float64(data)
			}
		case float64:
			if min > data {
				min = float64(data)
			}
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
		}
	}

	return min
}

func calculateMax(c *AggregateCondition) float64 {
	var max float64
	firstEvent := c.EventBuffer[c.start]
	switch data := firstEvent.Value.(type) {
	case int:
		max = float64(data)
	case float64:
		max = data
	default:
		panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
	}

	for i := range c.eventCount {
		event := c.EventBuffer[(c.start+i)%c.MaxSize]
		data := event.Value
		switch data := data.(type) {
		case int:
			if max < float64(data) {
				max = float64(data)
			}
		case float64:
			if max < data {
				max = float64(data)
			}
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
		}
	}

	return max
}

func calculateSum(c *AggregateCondition) float64 {
	var sum float64
	for i := range c.eventCount {
		event := c.EventBuffer[(c.start+i)%c.MaxSize]
		data := event.Value
		switch data := data.(type) {
		case int:
			sum += float64(data)
		case float64:
			sum += data
		default:
			panic(fmt.Sprintf("unsupported type: %s", reflect.TypeOf(data)))
		}
	}
	return sum
}
