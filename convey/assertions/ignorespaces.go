package assertions

import (
	"fmt"
	"reflect"
)

// ShouldBeEqualIgnoringSpaces receives exactly two string parameters and ensures that the order
// of non-space characters in the two strings is identical. Only one character ' '
// is considered a space. This is not a general white-space ignoring routine! Differences in tabs or
// newlines *will* be noticed, and will cause the two strings to look different.
func ShouldBeEqualIgnoringSpaces(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}

	value, valueIsString := actual.(string)
	expec, expecIsString := expected[0].(string)

	if !valueIsString || !expecIsString {
		return fmt.Sprintf(shouldBothBeStrings, reflect.TypeOf(actual), reflect.TypeOf(expected[0]))
	}

	if equalIgnoringSpaces(value, expec) {
		return success
	} else {
		return fmt.Sprintf(shouldHaveBeenEqualIgnoringSpaces, value, expec)
	}
}

func equalIgnoringSpaces(r, s string) bool {
	nextr := 0
	nexts := 0

	for {
		// skip past spaces in both r and s
		for nextr < len(r) {
			if r[nextr] == ' ' {
				nextr++
			} else {
				break
			}
		}

		for nexts < len(s) {
			if s[nexts] == ' ' {
				nexts++
			} else {
				break
			}
		}

		if nextr >= len(r) && nexts >= len(s) {
			return true
		}

		if nextr >= len(r) {
			return false
		}
		if nexts >= len(s) {
			return false
		}

		if r[nextr] != s[nexts] {
			return false
		}
		nextr++
		nexts++
	}

	return false
}
