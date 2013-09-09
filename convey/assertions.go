package convey

import (
	"fmt"
	"reflect"
)

import (
	"github.com/jacobsa/oglematchers"
)

// assertion is an alias for a function with a signature that the convey.So()
// method can handle. Any future or custom assertions should conform to this
// method signature. The return value should be an empty string if the assertion
// passes and a well-formed failure message if not.
type assertion func(actual interface{}, expected ...interface{}) string

// ShouldEqual receives exactly two parameters and does a simple equality (==) check.
func ShouldEqual(actual interface{}, expected ...interface{}) (message string) {
	if message = onlyOne(expected); message != success {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			message = fmt.Sprintf(shouldHaveBeenEqual, actual, expected[0])
		}
	}()

	matcher := oglematchers.Equals(expected[0])
	matchError := matcher.Matches(actual)
	if matchError != nil {
		message = fmt.Sprintf(shouldHaveBeenEqual, actual, expected[0])
		return
	}

	return
}

// ShouldNotEqual receives exactly two parameters and does a simple inequality (!=) check.
func ShouldNotEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	} else if ShouldEqual(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveBeenEqual, actual, expected[0])
	}
	return success
}

// ShouldResemble receives exactly two parameters and does a deep equal check (see reflect.DeepEqual)
func ShouldResemble(actual interface{}, expected ...interface{}) (message string) {
	if message = onlyOne(expected); message != success {
		return
	}

	matcher := oglematchers.DeepEquals(expected[0])
	matchError := matcher.Matches(actual)
	if matchError != nil {
		message = fmt.Sprintf(shouldHaveResembled, actual, expected[0])
		return
	}
	
	return
}

// ShouldNotResemble receives exactly two parameters and does an inverse deep equal check (see reflect.DeepEqual)
func ShouldNotResemble(actual interface{}, expected ...interface{}) (message string) {
	if message = onlyOne(expected); message != success {
		return
	} else if ShouldResemble(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveResembled, actual, expected[0])
	}
	return success
}

// ShouldBeNil receives a single parameter and does a nil check.
func ShouldBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if actual == nil {
		return success
	} else if interfaceIsNilPointer(actual) {
		return success
	}
	return fmt.Sprintf(shouldHaveBeenNil, actual)
}
func interfaceIsNilPointer(actual interface{}) bool {
	value := reflect.ValueOf(actual)
	return value.Kind() == reflect.Ptr && value.Pointer() == 0
}

// ShouldNotBeNil receives a single parameter and ensures it is not nil.
func ShouldNotBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if ShouldBeNil(actual) == success {
		return fmt.Sprintf(shouldNotHaveBeenNil, actual)
	}
	return success
}

// ShouldBeTrue receives a single parameter and ensures it is true.
func ShouldBeTrue(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if actual != true {
		return fmt.Sprintf(shouldHaveBeenTrue, actual)
	}
	return success
}

// ShouldBeFalse receives a single parameter and ensures it is false.
func ShouldBeFalse(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if actual != false {
		return fmt.Sprintf(shouldHaveBeenFalse, actual)
	}
	return success
}
