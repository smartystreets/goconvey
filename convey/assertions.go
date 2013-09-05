package convey

import (
	"fmt"
	_ "reflect"
)

// assertion is an alias for a function with a signature that the convey.So()
// method can handle. Any future or custom assertions should conform to this
// method signature. The return value should be an empty string if the assertion
// passes and a well-formed failure message if not.
type assertion func(actual interface{}, expected ...interface{}) string

// ShouldEqual receives exactly two parameters and does a simple equality (==) check.
func ShouldEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != "" {
		return fail
	} else if actual != expected[0] {
		return fmt.Sprintf(shouldHaveBeenEqual, actual, expected[0])
	}
	return success
}

// ShouldNotEqual receives exactly two parameters and does a simple inequality (!=) check.
func ShouldNotEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != "" {
		return fail
	} else if actual == expected[0] {
		return fmt.Sprintf(shouldNotHaveBeenEqual, actual, expected[0])
	}
	return success
}

// ShouldBeNil receives a single parameter and does a nil check.
func ShouldBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != "" {
		return fail
	} else if actual != nil {
		return fmt.Sprintf(shouldHaveBeenNil, actual)
	}
	return success
}

// ShouldNotBeNil receives a single parameter and ensures it is not nil.
func ShouldNotBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != "" {
		return fail
	} else if actual == nil {
		return fmt.Sprintf(shouldNotHaveBeenNil, actual)
	}
	return success
}

// ShouldBeTrue receives a single parameter and ensures it is true.
func ShouldBeTrue(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != "" {
		return fail
	} else if actual != true {
		return fmt.Sprintf(shouldHaveBeenTrue, actual)
	}
	return success
}

// ShouldBeFalse receives a single parameter and ensures it is false.
func ShouldBeFalse(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != "" {
		return fail
	} else if actual != false {
		return fmt.Sprintf(shouldHaveBeenFalse, actual)
	}
	return success
}
