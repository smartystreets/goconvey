package convey

import (
	_ "fmt"
	_ "reflect"
	"runtime"
	"testing"
)

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

	pass(t, so(thing1{}, ShouldEqual, thing1{}))
	fail(t, so(&thing1{"hi"}, ShouldEqual, &thing1{"hi"}), "'&{hi}' should equal '&{hi}' (but it doesn't)!")

	fail(t, so(thing1{}, ShouldEqual, thing2{}), "'{}' should equal '{}' (but it doesn't)!")
}

func TestShouldNotEqual(t *testing.T) {
	fail(t, so(1, ShouldNotEqual), "This expectation requires at least one comparison value (none provided).")
	fail(t, so(1, ShouldNotEqual, 1, 2), "This expectation only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldNotEqual, 1, 2, 3), "This expectation only accepts 1 value to be compared (and 3 were provided).")

	pass(t, so(1, ShouldNotEqual, 2))
	fail(t, so(1, ShouldNotEqual, 1), "'1' should NOT equal '1' (but it does)!")

	pass(t, so(true, ShouldNotEqual, false))
	fail(t, so(true, ShouldNotEqual, true), "'true' should NOT equal 'true' (but it does)!")

	pass(t, so("hi", ShouldNotEqual, "bye"))
	fail(t, so("hi", ShouldNotEqual, "hi"), "'hi' should NOT equal 'hi' (but it does)!")

	pass(t, so(&thing1{"hi"}, ShouldNotEqual, &thing1{"hi"}))
	fail(t, so(thing1{}, ShouldNotEqual, thing1{}), "'{}' should NOT equal '{}' (but it does)!")

	pass(t, so(thing1{}, ShouldNotEqual, thing2{}))
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
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Expectation should have passed but failed (see line %d): '%s'", line, result)
	}
}

func fail(t *testing.T, actual string, expected string) {
	if actual != expected {
		if actual == "" {
			actual = "(empty)"
		}
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Expectation should have failed but passed (see line %d). \nExpected: %s\nActual:   %s\n",
			line, expected, actual)
	}
}

func so(actual interface{}, match expectation, expected ...interface{}) string {
	return match(actual, expected...)
}

type thing1 struct {
	a string
}
type thing2 struct {
	a string
}
