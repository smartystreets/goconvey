package convey

import (
	"testing"
)

func TestShouldEqual(t *testing.T) {
	fail(t, so(1, ShouldEqual), "This expectation requires a second value (only one provided: '1').")
	fail(t, so(1, ShouldEqual, 1, 2), "This expectation only accepts 2 values to be compared (and 3 were provided).")
	fail(t, so(1, ShouldEqual, 1, 2, 3), "This expectation only accepts 2 values to be compared (and 4 were provided).")

	pass(t, so(1, ShouldEqual, 1))
	fail(t, so(1, ShouldEqual, 2), "'1' should equal '2' (but it doesn't)!")

	pass(t, so(true, ShouldEqual, true))
	fail(t, so(true, ShouldEqual, false), "'true' should equal 'false' (but it doesn't)!")

	pass(t, so("hi", ShouldEqual, "hi"))
	fail(t, so("hi", ShouldEqual, "bye"), "'hi' should equal 'bye' (but it doesn't)!")
}

func pass(t *testing.T, result string) {
	const PASS = ""
	if result != PASS {
		t.Errorf("Expectation should have passed but failed: '%s'", result)
	}
}

func fail(t *testing.T, actual string, expected string) {
	if actual != expected {
		if actual == "" {
			actual = "(empty)"
		}
		t.Errorf("Expectation should have failed. \nExpected: %s\nActual:   %s\n",
			expected, actual)
	}
}

func so(actual interface{}, match expectation, expected ...interface{}) string {
	return match(actual, expected)
}
