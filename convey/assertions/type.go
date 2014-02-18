package assertions

import (
	"fmt"
	"reflect"
)

// ShouldHaveSameTypeAs receives exactly two parameters and compares their underlying types for equality.
func ShouldHaveSameTypeAs(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	first := reflect.TypeOf(actual)
	second := reflect.TypeOf(expected[0])

	if equal := ShouldEqual(first, second); equal != success {
		return serializer.serialize(second, first, fmt.Sprintf(shouldHaveBeenA, actual, second, first))
	}
	return success
}

// ShouldNotHaveSameTypeAs receives exactly two parameters and compares their underlying types for inequality.
func ShouldNotHaveSameTypeAs(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	first := reflect.TypeOf(actual)
	second := reflect.TypeOf(expected[0])

	if equal := ShouldEqual(first, second); equal == success {
		return fmt.Sprintf(shouldNotHaveBeenA, actual, second)
	}
	return success
}
