package convey

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestShouldEqual(t *testing.T) {
	fail(t, so(1, ShouldEqual), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(1, ShouldEqual, 1, 2), "This expectation only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldEqual, 1, 2, 3), "This expectation only accepts 1 value to be compared (and 3 were provided).")

	pass(t, so(1, ShouldEqual, 1))
	fail(t, so(1, ShouldEqual, 2), "Expected '1' to equal '2' (but it didn't)!")

	pass(t, so(true, ShouldEqual, true))
	fail(t, so(true, ShouldEqual, false), "Expected 'true' to equal 'false' (but it didn't)!")

	pass(t, so("hi", ShouldEqual, "hi"))
	fail(t, so("hi", ShouldEqual, "bye"), "Expected 'hi' to equal 'bye' (but it didn't)!")

	pass(t, so(42, ShouldEqual, uint(42)))

	fail(t, so(Thing1{}, ShouldEqual, Thing1{}), "Expected '{}' to equal '{}' (but it didn't)!")
	fail(t, so(Thing1{"hi"}, ShouldEqual, Thing1{"hi"}), "Expected '{hi}' to equal '{hi}' (but it didn't)!")
	fail(t, so(&Thing1{"hi"}, ShouldEqual, &Thing1{"hi"}), "Expected '&{hi}' to equal '&{hi}' (but it didn't)!")

	fail(t, so(Thing1{}, ShouldEqual, Thing2{}), "Expected '{}' to equal '{}' (but it didn't)!")
}

func TestShouldNotEqual(t *testing.T) {
	fail(t, so(1, ShouldNotEqual), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(1, ShouldNotEqual, 1, 2), "This expectation only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldNotEqual, 1, 2, 3), "This expectation only accepts 1 value to be compared (and 3 were provided).")

	pass(t, so(1, ShouldNotEqual, 2))
	fail(t, so(1, ShouldNotEqual, 1), "Expected '1' to NOT equal '1' (but it did)!")

	pass(t, so(true, ShouldNotEqual, false))
	fail(t, so(true, ShouldNotEqual, true), "Expected 'true' to NOT equal 'true' (but it did)!")

	pass(t, so("hi", ShouldNotEqual, "bye"))
	fail(t, so("hi", ShouldNotEqual, "hi"), "Expected 'hi' to NOT equal 'hi' (but it did)!")

	pass(t, so(&Thing1{"hi"}, ShouldNotEqual, &Thing1{"hi"}))
	pass(t, so(Thing1{"hi"}, ShouldNotEqual, Thing1{"hi"}))
	pass(t, so(Thing1{}, ShouldNotEqual, Thing1{}))
	pass(t, so(Thing1{}, ShouldNotEqual, Thing2{}))
}

func TestShouldResemble(t *testing.T) {
	fail(t, so(Thing1{"hi"}, ShouldResemble), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}, Thing1{"hi"}), "This expectation only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}))
	fail(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"bye"}), "Expected '{hi}' to resemble '{bye}' (but it didn't)!")
}

func TestShouldNotResemble(t *testing.T) {
	fail(t, so(Thing1{"hi"}, ShouldNotResemble), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}, Thing1{"hi"}), "This expectation only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"bye"}))
	fail(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}), "Expected '{hi}' to NOT resemble '{hi}' (but it did)!")
}

func TestShouldPointTo(t *testing.T) {
	t1 := &Thing1{}
	t2 := t1
	t3 := &Thing1{}

	pointer1 := reflect.ValueOf(t1).Pointer()
	pointer3 := reflect.ValueOf(t3).Pointer()

	fail(t, so(t1, ShouldPointTo), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(t1, ShouldPointTo, t2, t3), "This expectation only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(t1, ShouldPointTo, t2))
	fail(t, so(t1, ShouldPointTo, t3), fmt.Sprintf("Expected '&{}' (address: '%v') and '&{}' (address: '%v') to be the same address (but their weren't)!", pointer1, pointer3))

	t4 := Thing1{}
	t5 := t4

	fail(t, so(t4, ShouldPointTo, t5), "Both arguments should be pointers (the first was not)!")
	fail(t, so(&t4, ShouldPointTo, t5), "Both arguments should be pointers (the second was not)!")
	fail(t, so(nil, ShouldPointTo, nil), "Both arguments should be pointers (the first was nil)!")
	fail(t, so(&t4, ShouldPointTo, nil), "Both arguments should be pointers (the second was nil)!")
}

func TestShouldNotPointTo(t *testing.T) {
	t1 := &Thing1{}
	t2 := t1
	t3 := &Thing1{}

	pointer1 := reflect.ValueOf(t1).Pointer()

	fail(t, so(t1, ShouldNotPointTo), "This expectation requires exactly one comparison value (none provided).")
	fail(t, so(t1, ShouldNotPointTo, t2, t3), "This expectation only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(t1, ShouldNotPointTo, t3))
	fail(t, so(t1, ShouldNotPointTo, t2), fmt.Sprintf("Expected '&{}' and '&{}' to be different references (but they matched: '%v')!", pointer1))

	t4 := Thing1{}
	t5 := t4

	fail(t, so(t4, ShouldNotPointTo, t5), "Both arguments should be pointers (the first was not)!")
	fail(t, so(&t4, ShouldNotPointTo, t5), "Both arguments should be pointers (the second was not)!")
	fail(t, so(nil, ShouldNotPointTo, nil), "Both arguments should be pointers (the first was nil)!")
	fail(t, so(&t4, ShouldNotPointTo, nil), "Both arguments should be pointers (the second was nil)!")
}

func TestShouldBeNil(t *testing.T) {
	fail(t, so(nil, ShouldBeNil, nil, nil, nil), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldBeNil, nil), "This expectation does not allow for user-supplied comparison values.")

	pass(t, so(nil, ShouldBeNil))
	fail(t, so(1, ShouldBeNil), "Expected '1' to be nil (but it wasn't)!")

	var thing Thinger
	pass(t, so(thing, ShouldBeNil))
	thing = &Thing{}
	fail(t, so(thing, ShouldBeNil), "Expected '&{}' to be nil (but it wasn't)!")

	var thingOne *Thing1
	pass(t, so(thingOne, ShouldBeNil))
}

func TestShouldNotBeNil(t *testing.T) {
	fail(t, so(nil, ShouldNotBeNil, nil, nil, nil), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldNotBeNil, nil), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(nil, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	pass(t, so(1, ShouldNotBeNil))

	var thing Thinger
	fail(t, so(thing, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	thing = &Thing{}
	pass(t, so(thing, ShouldNotBeNil))
}

func TestShouldBeTrue(t *testing.T) {
	fail(t, so(true, ShouldBeTrue, 1, 2, 3), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(true, ShouldBeTrue, 1), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(false, ShouldBeTrue), "Expected 'true' (not 'false')!")
	fail(t, so(1, ShouldBeTrue), "Expected 'true' (not '1')!")
	pass(t, so(true, ShouldBeTrue))
}

func TestShouldBeFalse(t *testing.T) {
	fail(t, so(false, ShouldBeFalse, 1, 2, 3), "This expectation does not allow for user-supplied comparison values.")
	fail(t, so(false, ShouldBeFalse, 1), "This expectation does not allow for user-supplied comparison values.")

	fail(t, so(true, ShouldBeFalse), "Expected 'false' (not 'true')!")
	fail(t, so(1, ShouldBeFalse), "Expected 'false' (not '1')!")
	pass(t, so(false, ShouldBeFalse))
}

func pass(t *testing.T, result string) {
	if result != success {
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
