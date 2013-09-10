package assertions

import "fmt"

const (
	success             = ""
	needOneValue        = "This assertion requires exactly one comparison value (none provided)."
	onlyAcceptsOneValue = "This assertion only accepts 1 value to be compared (and %v were provided)."
	noValuesAccepted    = "This assertion does not allow for user-supplied comparison values."
	needTwoValues       = "This assertion requires exactly 2 comparison values (you provided %d)."
)

func onlyTwo(expected []interface{}) string {
	if len(expected) != 2 {
		return fmt.Sprintf(needTwoValues, len(expected))
	}
	return success
}

func onlyOne(expected []interface{}) string {
	switch {
	case len(expected) == 0:
		return needOneValue
	case len(expected) > 1:
		return fmt.Sprintf(onlyAcceptsOneValue, len(expected))
	}
	return success
}

func none(expected []interface{}) string {
	if len(expected) > 0 {
		return noValuesAccepted
	}
	return success
}
