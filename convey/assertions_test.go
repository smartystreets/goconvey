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

	pass(t, so(Thing1{}, ShouldEqual, Thing1{}))
	pass(t, so(Thing1{"hi"}, ShouldEqual, Thing1{"hi"}))
	fail(t, so(&Thing1{"hi"}, ShouldEqual, &Thing1{"hi"}), "'&{hi}' should equal '&{hi}' (but it doesn't)!")

	fail(t, so(Thing1{}, ShouldEqual, Thing2{}), "'{}' should equal '{}' (but it doesn't)!")
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

	pass(t, so(&Thing1{"hi"}, ShouldNotEqual, &Thing1{"hi"}))
	fail(t, so(Thing1{"hi"}, ShouldNotEqual, Thing1{"hi"}), "'{hi}' should NOT equal '{hi}' (but it does)!")
	fail(t, so(Thing1{}, ShouldNotEqual, Thing1{}), "'{}' should NOT equal '{}' (but it does)!")

	pass(t, so(Thing1{}, ShouldNotEqual, Thing2{}))
}

func TestShouldBeNil(t *testing.T) {
	fail(t, so(nil, ShouldBeNil, nil, nil, nil), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldBeNil, nil), "This expectation does not allow for user-supplied comparison values.")

	pass(t, so(nil, ShouldBeNil))
	fail(t, so(1, ShouldBeNil), "'1' should have been nil (but it wasn't)!")

	var thing Thinger
	pass(t, so(thing, ShouldBeNil))
	thing = &Thing{}
	fail(t, so(thing, ShouldBeNil), "'&{}' should have been nil (but it wasn't)!")
}

func TestShouldNotBeNil(t *testing.T) {
	fail(t, so(nil, ShouldNotBeNil, nil, nil, nil), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldNotBeNil, nil), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(nil, ShouldNotBeNil), "'<nil>' should NOT have been nil (but it was)!")
	pass(t, so(1, ShouldNotBeNil))

	var thing Thinger
	fail(t, so(thing, ShouldNotBeNil), "'<nil>' should NOT have been nil (but it was)!")
	thing = &Thing{}
	pass(t, so(thing, ShouldNotBeNil))
}

func TestShouldBeTrue(t *testing.T) {
	fail(t, so(true, ShouldBeTrue, 1, 2, 3), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(true, ShouldBeTrue, 1), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(false, ShouldBeTrue), "Should have been 'true', not 'false'.")
	fail(t, so(1, ShouldBeTrue), "Should have been 'true', not '1'.")
	pass(t, so(true, ShouldBeTrue))
}

func TestShouldBeFalse(t *testing.T) {
	fail(t, so(false, ShouldBeFalse, 1, 2, 3), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(false, ShouldBeFalse, 1), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(true, ShouldBeFalse), "Should have been 'false', not 'true'.")
	fail(t, so(1, ShouldBeFalse), "Should have been 'false', not '1'.")
	pass(t, so(false, ShouldBeFalse))
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

func so(actual interface{}, assert assertion, expected ...interface{}) string {
	return assert(actual, expected...)
}

type Thing1 struct {
	a string
}
type Thing2 struct {
	a string
}

type Thinger interface {
	Hi()
}

type Thing struct{}

func (self *Thing) Hi() {}
