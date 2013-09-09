package convey

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/jacobsa/oglematchers"
)

// assertion is an alias for a function with a signature that the convey.So()
// method can handle. Any future or custom assertions should conform to this
// method signature. The return value should be an empty string if the assertion
// passes and a well-formed failure message if not.
type assertion func(actual interface{}, expected ...interface{}) string

// ShouldEqual receives exactly two parameters and does an equality check.
func ShouldEqual(actual interface{}, expected ...interface{}) string {
	if message := onlyOne(expected); message != success {
		return message
	}
	return shouldEqual(actual, expected[0])
}
func shouldEqual(actual, expected interface{}) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = fmt.Sprintf(shouldHaveBeenEqual, actual, expected)
		}
	}()

	if matchError := oglematchers.Equals(expected).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenEqual, actual, expected)
	}

	return success
}

// ShouldNotEqual receives exactly two parameters and does an inequality check.
func ShouldNotEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != success {
		return fail
	} else if ShouldEqual(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveBeenEqual, actual, expected[0])
	}
	return success
}

// ShouldResemble receives exactly two parameters and does a deep equal check (see reflect.DeepEqual)
func ShouldResemble(actual interface{}, expected ...interface{}) string {
	if message := onlyOne(expected); message != success {
		return message
	}

	matcher := oglematchers.DeepEquals(expected[0])
	matchError := matcher.Matches(actual)
	if matchError != nil {
		return fmt.Sprintf(shouldHaveResembled, actual, expected[0])
	}

	return success
}

// ShouldNotResemble receives exactly two parameters and does an inverse deep equal check (see reflect.DeepEqual)
func ShouldNotResemble(actual interface{}, expected ...interface{}) string {
	if message := onlyOne(expected); message != success {
		return message
	} else if ShouldResemble(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveResembled, actual, expected[0])
	}
	return success
}

// ShouldPointTo receives exactly two parameters and checks to see that they point to the same address.
func ShouldPointTo(actual interface{}, expected ...interface{}) string {
	if message := onlyOne(expected); message != success {
		return message
	}
	return shouldPointTo(actual, expected[0])

}
func shouldPointTo(actual, expected interface{}) string {
	actualValue := reflect.ValueOf(actual)
	expectedValue := reflect.ValueOf(expected)

	if ShouldNotBeNil(actual) != success {
		return fmt.Sprintf(shouldHaveBeenNonNilPointer, "first", "nil")
	} else if ShouldNotBeNil(expected) != success {
		return fmt.Sprintf(shouldHaveBeenNonNilPointer, "second", "nil")
	} else if actualValue.Kind() != reflect.Ptr {
		return fmt.Sprintf(shouldHaveBeenNonNilPointer, "first", "not")
	} else if expectedValue.Kind() != reflect.Ptr {
		return fmt.Sprintf(shouldHaveBeenNonNilPointer, "second", "not")
	} else if ShouldEqual(actualValue.Pointer(), expectedValue.Pointer()) != success {
		return fmt.Sprintf(shouldHavePointedTo,
			actual, reflect.ValueOf(actual).Pointer(),
			expected, reflect.ValueOf(expected).Pointer())
	}
	return success
}

// ShouldNotPointTo receives exactly two parameters and checks to see that they point to different addresess.
func ShouldNotPointTo(actual interface{}, expected ...interface{}) string {
	if message := onlyOne(expected); message != success {
		return message
	}
	compare := ShouldPointTo(actual, expected[0])
	if strings.HasPrefix(compare, shouldBePointers) {
		return compare
	} else if compare == success {
		return fmt.Sprintf(shouldNotHavePointedTo, actual, expected[0], reflect.ValueOf(actual).Pointer())
	}
	return success
}

// ShouldBeNil receives a single parameter and ensures that it is nil.
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

// ShouldNotBeNil receives a single parameter and ensures that it is not nil.
func ShouldNotBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if ShouldBeNil(actual) == success {
		return fmt.Sprintf(shouldNotHaveBeenNil, actual)
	}
	return success
}

// ShouldBeTrue receives a single parameter and ensures that it is true.
func ShouldBeTrue(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if actual != true {
		return fmt.Sprintf(shouldHaveBeenTrue, actual)
	}
	return success
}

// ShouldBeFalse receives a single parameter and ensures that it is false.
func ShouldBeFalse(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != success {
		return fail
	} else if actual != false {
		return fmt.Sprintf(shouldHaveBeenFalse, actual)
	}
	return success
}

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
