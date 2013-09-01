package convey

import "testing"

func TestShouldEqual(t *testing.T) {
	fail(t, so(1, ShouldEqual), "This expectation requires at least one comparison value (none provided).")
	fail(t, so(1, ShouldEqual, 1, 2), "This expectation only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldEqual, 1, 2, 3), "This expectation only accepts 1 value to be compared (and 3 were provided).")

	pass(t, so(1, ShouldEqual, 1))
	fail(t, so(1, ShouldEqual, 2), "'1' should equal '2' (but it doesn't)!")

	pass(t, so(true, ShouldEqual, true))
	fail(t, so(true, ShouldEqual, false), "'true' should equal 'false' (but it doesn't)!")

	pass(t, so("hi", ShouldEqual, "hi"))
	fail(t, so("hi", ShouldEqual, "bye"), "'hi' should equal 'bye' (but it doesn't)!")

	// TODO: compare structs
}

func TestShouldBeNil(t *testing.T) {
	fail(t, so(nil, ShouldBeNil, nil, nil, nil), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldBeNil, nil), "This expectation does not allow for user-supplied comparison values.")

	pass(t, so(nil, ShouldBeNil))
	fail(t, so(1, ShouldBeNil), "'1' should have been nil!")
}

func pass(t *testing.T, result string) {
	const PASS = success
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
