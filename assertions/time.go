package assertions

import (
	"fmt"
	"time"
)

func ShouldHappenBefore(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}

	if !actualTime.Before(expectedTime) {
		return fmt.Sprintf(shouldHaveHappenedBefore, actualTime, expectedTime, actualTime.Sub(expectedTime))
	}

	return success
}

func ShouldHappenOnOrBefore(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}

func ShouldHappenAfter(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}

func ShouldHappenOnOrAfter(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}

func ShouldHappenBetween(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}

func ShouldHappenOnOrBetween(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}

func ShouldHappenWithin(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	actualTime, firstOk := actual.(time.Time)
	expectedTime, secondOk := expected[0].(time.Time)

	if !firstOk || !secondOk {
		return shouldUseTimes
	}
	return fmt.Sprintf("%v %v", actualTime, expectedTime) // TODO
}
