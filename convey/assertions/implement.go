package assertions

import (
	"fmt"
	"reflect"
)

// ShouldImplement receives exactly two parameters and compares their underlying types for equality.
func ShouldImplement(actual interface{}, expectedList ...interface{}) string {
	if fail := need(1, expectedList); fail != success {
		return fail
	}
	if fail := ShouldNotBeNil(expectedList); fail != success {
		return fail
	}
	expected := expectedList[0]
	if fail := ShouldBeNil(expected); fail != success {
		return "This assertion requires a pointer with the interface type"
	}
	expectedType := reflect.TypeOf(expected)
	if fail := ShouldNotBeNil(expectedType); fail != success {
		return "This assertion requires a pointer with the interface type"
	}
	expectedInterface := expectedType.Elem()
	actualType := reflect.TypeOf(actual)
	if actualType == nil{
		return fmt.Sprintf(shouldImplement, actual, expectedInterface, actualType)
	}
	if fail := ShouldEqual(actualType.Kind(), reflect.Ptr); fail != success {
		return fmt.Sprintf(shouldImplement, actual, expectedInterface, actualType)
	}
	if actualType == nil {
		return fmt.Sprintf(shouldImplement, actual, expectedInterface, actualType)
	}
	if !actualType.Implements(expectedInterface) {
		return fmt.Sprintf(shouldImplement, actualType, expectedInterface, actualType)
	}
	return success
}
