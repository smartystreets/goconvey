package convey

import "fmt"

const success = ""
const needAtLeastOneValue = "This expectation requires at least one comparison value (none provided)."
const onlyAcceptsOneValue = "This expectation only accepts 1 value to be compared (and %v were provided)."
const noValuesAccepted = "This expectation does not allow for user-supplied comparison values."
const shouldHaveBeenEqual = "'%v' should equal '%v' (but it doesn't)!"
const shouldNotHaveBeenEqual = "'%v' should NOT equal '%v' (but it does)!"
const shouldHaveBeenNil = "'%v' should have been nil!"

func onlyOne(expected []interface{}) string {
	switch {
	case len(expected) == 0:
		return needAtLeastOneValue
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
