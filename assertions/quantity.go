package assertions

import (
	"fmt"
)

import (
	"github.com/smartystreets/oglematchers"
)

// ShouldBeGreaterThan receives exactly two parameters and ensures that the first is greater than the second.
func ShouldBeGreaterThan(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	}

	if matchError := oglematchers.GreaterThan(expected[0]).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenGreater, actual, expected[0])
	}
	return success
}

// ShouldBeGreaterThanOrEqualTo receives exactly two parameters and ensures that the first is greater than or equal to the second.
func ShouldBeGreaterThanOrEqualTo(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	} else if matchError := oglematchers.GreaterOrEqual(expected[0]).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenGreaterOrEqual, actual, expected[0])
	}
	return success
}

// ShouldBeLessThan receives exactly two parameters and ensures that the first is less than the second.
func ShouldBeLessThan(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	} else if matchError := oglematchers.LessThan(expected[0]).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenLess, actual, expected[0])
	}
	return success
}

// ShouldBeLessThan receives exactly two parameters and ensures that the first is less than or equal to the second.
func ShouldBeLessThanOrEqualTo(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	} else if matchError := oglematchers.LessOrEqual(expected[0]).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenLess, actual, expected[0])
	}
	return success
}

// ShouldBeBetween receives exactly three parameters: an actual value, a lower bound, and an upper bound.
// It ensures that the actual value is between both bounds or is equal to one of the bounds.
func ShouldBeBetween(actual interface{}, expected ...interface{}) string {
	if fail := onlyTwo(expected); fail != success {
		return fail
	}
	return "asdf"
}
