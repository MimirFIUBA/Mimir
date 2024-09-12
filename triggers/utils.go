package triggers

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
