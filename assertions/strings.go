package assertions

import (
	"fmt"
	"reflect"
	"strings"
)

// ShouldStartWith receives exactly 2 string parameters and ensures that the first starts with the second.
func ShouldStartWith(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	value, valueIsString := actual.(string)
	prefix, prefixIsString := expected[0].(string)

	if !valueIsString || !prefixIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	return shouldContain(value, prefix)
}
func shouldContain(value, prefix string) string {
	if !strings.HasPrefix(value, prefix) {
		return fmt.Sprintf(shouldHaveStartedWith, value, prefix)
	}
	return success
}

// ShouldNotStartWith receives exactly 2 string parameters and ensures that the first does not start with the second.
func ShouldNotStartWith(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	value, valueIsString := actual.(string)
	prefix, prefixIsString := expected[0].(string)

	if !valueIsString || !prefixIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	return shouldNotContain(value, prefix)
}
func shouldNotContain(value, prefix string) string {
	if strings.Contains(value, prefix) {
		if value == "" {
			value = "<empty>"
		}
		if prefix == "" {
			prefix = "<empty>"
		}
		return fmt.Sprintf(shouldNotHaveStartedWith, value, prefix)
	}
	return success
}
