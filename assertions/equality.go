package assertions

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	"github.com/jacobsa/oglematchers"
)

// ShouldEqual receives exactly two parameters and does an equality check.
func ShouldEqual(actual interface{}, expected ...interface{}) string {
	if message := need(1, expected); message != success {
		return message
	}
	return shouldEqual(actual, expected[0])
}
func shouldEqual(actual, expected interface{}) (message string) {
	defer func() {
		if r := recover(); r != nil {
			message = fmt.Sprintf(shouldHaveBeenEqual, expected, actual)
		}
	}()

	if matchError := oglematchers.Equals(expected).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveBeenEqual, expected, actual)
	}

	return success
}

// ShouldNotEqual receives exactly two parameters and does an inequality check.
func ShouldNotEqual(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	} else if ShouldEqual(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveBeenEqual, actual, expected[0])
	}
	return success
}

// ShouldResemble receives exactly two parameters and does a deep equal check (see reflect.DeepEqual)
func ShouldResemble(actual interface{}, expected ...interface{}) string {
	if message := need(1, expected); message != success {
		return message
	}

	if matchError := oglematchers.DeepEquals(expected[0]).Matches(actual); matchError != nil {
		return fmt.Sprintf(shouldHaveResembled, expected[0], actual)
	}

	return success
}

// ShouldNotResemble receives exactly two parameters and does an inverse deep equal check (see reflect.DeepEqual)
func ShouldNotResemble(actual interface{}, expected ...interface{}) string {
	if message := need(1, expected); message != success {
		return message
	} else if ShouldResemble(actual, expected[0]) == success {
		return fmt.Sprintf(shouldNotHaveResembled, actual, expected[0])
	}
	return success
}

// ShouldPointTo receives exactly two parameters and checks to see that they point to the same address.
func ShouldPointTo(actual interface{}, expected ...interface{}) string {
	if message := need(1, expected); message != success {
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
	if message := need(1, expected); message != success {
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
	if fail := need(0, expected); fail != success {
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
	if fail := need(0, expected); fail != success {
		return fail
	} else if ShouldBeNil(actual) == success {
		return fmt.Sprintf(shouldNotHaveBeenNil, actual)
	}
	return success
}

// ShouldBeTrue receives a single parameter and ensures that it is true.
func ShouldBeTrue(actual interface{}, expected ...interface{}) string {
	if fail := need(0, expected); fail != success {
		return fail
	} else if actual != true {
		return fmt.Sprintf(shouldHaveBeenTrue, actual)
	}
	return success
}

// ShouldBeFalse receives a single parameter and ensures that it is false.
func ShouldBeFalse(actual interface{}, expected ...interface{}) string {
	if fail := need(0, expected); fail != success {
		return fail
	} else if actual != false {
		return fmt.Sprintf(shouldHaveBeenFalse, actual)
	}
	return success
}
