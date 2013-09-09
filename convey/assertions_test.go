package convey

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
)

func TestShouldEqual(t *testing.T) {
	fail(t, so(1, ShouldEqual), needOneValue)
	fail(t, so(1, ShouldEqual, 1, 2), "This assertion only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldEqual, 1, 2, 3), "This assertion only accepts 1 value to be compared (and 3 were provided).")

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
	fail(t, so(1, ShouldNotEqual), needOneValue)
	fail(t, so(1, ShouldNotEqual, 1, 2), "This assertion only accepts 1 value to be compared (and 2 were provided).")
	fail(t, so(1, ShouldNotEqual, 1, 2, 3), "This assertion only accepts 1 value to be compared (and 3 were provided).")

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
	fail(t, so(Thing1{"hi"}, ShouldResemble), needOneValue)
	fail(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}, Thing1{"hi"}), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}))
	fail(t, so(Thing1{"hi"}, ShouldResemble, Thing1{"bye"}), "Expected '{hi}' to resemble '{bye}' (but it didn't)!")
}

func TestShouldNotResemble(t *testing.T) {
	fail(t, so(Thing1{"hi"}, ShouldNotResemble), needOneValue)
	fail(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}, Thing1{"hi"}), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"bye"}))
	fail(t, so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}), "Expected '{hi}' to NOT resemble '{hi}' (but it did)!")
}

func TestShouldPointTo(t *testing.T) {
	t1 := &Thing1{}
	t2 := t1
	t3 := &Thing1{}

	pointer1 := reflect.ValueOf(t1).Pointer()
	pointer3 := reflect.ValueOf(t3).Pointer()

	fail(t, so(t1, ShouldPointTo), needOneValue)
	fail(t, so(t1, ShouldPointTo, t2, t3), "This assertion only accepts 1 value to be compared (and 2 were provided).")

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

	fail(t, so(t1, ShouldNotPointTo), needOneValue)
	fail(t, so(t1, ShouldNotPointTo, t2, t3), "This assertion only accepts 1 value to be compared (and 2 were provided).")

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
	fail(t, so(nil, ShouldBeNil, nil, nil, nil), "This assertion does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldBeNil, nil), "This assertion does not allow for user-supplied comparison values.")

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
	fail(t, so(nil, ShouldNotBeNil, nil, nil, nil), "This assertion does not allow for user-supplied comparison values.")
	fail(t, so(nil, ShouldNotBeNil, nil), "This assertion does not allow for user-supplied comparison values.")

	fail(t, so(nil, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	pass(t, so(1, ShouldNotBeNil))

	var thing Thinger
	fail(t, so(thing, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	thing = &Thing{}
	pass(t, so(thing, ShouldNotBeNil))
}

func TestShouldBeTrue(t *testing.T) {
	fail(t, so(true, ShouldBeTrue, 1, 2, 3), "This assertion does not allow for user-supplied comparison values.")
	fail(t, so(true, ShouldBeTrue, 1), "This assertion does not allow for user-supplied comparison values.")

	fail(t, so(false, ShouldBeTrue), "Expected 'true' (not 'false')!")
	fail(t, so(1, ShouldBeTrue), "Expected 'true' (not '1')!")
	pass(t, so(true, ShouldBeTrue))
}

func TestShouldBeFalse(t *testing.T) {
	fail(t, so(false, ShouldBeFalse, 1, 2, 3), "This assertion does not allow for user-supplied comparison values.")
	fail(t, so(false, ShouldBeFalse, 1), "This assertion does not allow for user-supplied comparison values.")

	fail(t, so(true, ShouldBeFalse), "Expected 'false' (not 'true')!")
	fail(t, so(1, ShouldBeFalse), "Expected 'false' (not '1')!")
	pass(t, so(false, ShouldBeFalse))
}

func TestShouldBeGreaterThan(t *testing.T) {
	fail(t, so(1, ShouldBeGreaterThan), needOneValue)
	fail(t, so(1, ShouldBeGreaterThan, 0, 0), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(1, ShouldBeGreaterThan, 0))
	pass(t, so(1.1, ShouldBeGreaterThan, 1))
	pass(t, so(1, ShouldBeGreaterThan, uint(0)))
	pass(t, so("b", ShouldBeGreaterThan, "a"))

	fail(t, so(0, ShouldBeGreaterThan, 1), "Expected '0' to be greater than '1' (but it wasn't)!")
	fail(t, so(1, ShouldBeGreaterThan, 1.1), "Expected '1' to be greater than '1.1' (but it wasn't)!")
	fail(t, so(uint(0), ShouldBeGreaterThan, 1.1), "Expected '0' to be greater than '1.1' (but it wasn't)!")
	fail(t, so("a", ShouldBeGreaterThan, "b"), "Expected 'a' to be greater than 'b' (but it wasn't)!")
}

func TestShouldBeGreaterThanOrEqual(t *testing.T) {
	fail(t, so(1, ShouldBeGreaterThanOrEqualTo), needOneValue)
	fail(t, so(1, ShouldBeGreaterThanOrEqualTo, 0, 0), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(1, ShouldBeGreaterThanOrEqualTo, 1))
	pass(t, so(1.1, ShouldBeGreaterThanOrEqualTo, 1.1))
	pass(t, so(1, ShouldBeGreaterThanOrEqualTo, uint(1)))
	pass(t, so("b", ShouldBeGreaterThanOrEqualTo, "b"))

	pass(t, so(1, ShouldBeGreaterThanOrEqualTo, 0))
	pass(t, so(1.1, ShouldBeGreaterThanOrEqualTo, 1))
	pass(t, so(1, ShouldBeGreaterThanOrEqualTo, uint(0)))
	pass(t, so("b", ShouldBeGreaterThanOrEqualTo, "a"))

	fail(t, so(0, ShouldBeGreaterThanOrEqualTo, 1), "Expected '0' to be greater than or equal to '1' (but it wasn't)!")
	fail(t, so(1, ShouldBeGreaterThanOrEqualTo, 1.1), "Expected '1' to be greater than or equal to '1.1' (but it wasn't)!")
	fail(t, so(uint(0), ShouldBeGreaterThanOrEqualTo, 1.1), "Expected '0' to be greater than or equal to '1.1' (but it wasn't)!")
	fail(t, so("a", ShouldBeGreaterThanOrEqualTo, "b"), "Expected 'a' to be greater than or equal to 'b' (but it wasn't)!")
}

func TestShouldBeLessThan(t *testing.T) {
	fail(t, so(1, ShouldBeLessThan), needOneValue)
	fail(t, so(1, ShouldBeLessThan, 0, 0), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(0, ShouldBeLessThan, 1))
	pass(t, so(1, ShouldBeLessThan, 1.1))
	pass(t, so(uint(0), ShouldBeLessThan, 1))
	pass(t, so("a", ShouldBeLessThan, "b"))

	fail(t, so(1, ShouldBeLessThan, 0), "Expected '1' to be less than '0' (but it wasn't)!")
	fail(t, so(1.1, ShouldBeLessThan, 1), "Expected '1.1' to be less than '1' (but it wasn't)!")
	fail(t, so(1.1, ShouldBeLessThan, uint(0)), "Expected '1.1' to be less than '0' (but it wasn't)!")
	fail(t, so("b", ShouldBeLessThan, "a"), "Expected 'b' to be less than 'a' (but it wasn't)!")
}

func TestShouldBeLessThanOrEqualTo(t *testing.T) {
	fail(t, so(1, ShouldBeLessThanOrEqualTo), needOneValue)
	fail(t, so(1, ShouldBeLessThanOrEqualTo, 0, 0), "This assertion only accepts 1 value to be compared (and 2 were provided).")

	pass(t, so(1, ShouldBeLessThanOrEqualTo, 1))
	pass(t, so(1.1, ShouldBeLessThanOrEqualTo, 1.1))
	pass(t, so(uint(1), ShouldBeLessThanOrEqualTo, 1))
	pass(t, so("b", ShouldBeLessThanOrEqualTo, "b"))

	pass(t, so(0, ShouldBeLessThanOrEqualTo, 1))
	pass(t, so(1, ShouldBeLessThanOrEqualTo, 1.1))
	pass(t, so(uint(0), ShouldBeLessThanOrEqualTo, 1))
	pass(t, so("a", ShouldBeLessThanOrEqualTo, "b"))

	fail(t, so(1, ShouldBeLessThanOrEqualTo, 0), "Expected '1' to be less than '0' (but it wasn't)!")
	fail(t, so(1.1, ShouldBeLessThanOrEqualTo, 1), "Expected '1.1' to be less than '1' (but it wasn't)!")
	fail(t, so(1.1, ShouldBeLessThanOrEqualTo, uint(0)), "Expected '1.1' to be less than '0' (but it wasn't)!")
	fail(t, so("b", ShouldBeLessThanOrEqualTo, "a"), "Expected 'b' to be less than 'a' (but it wasn't)!")
}

func TestShouldBeBetween(t *testing.T) {
	fail(t, so(1, ShouldBeBetween), "This assertion requires exactly 2 comparison values (you provided 0).")
	fail(t, so(1, ShouldBeBetween, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	pass(t, so(9, ShouldBeBetween, 8, 12))
	pass(t, so(11, ShouldBeBetween, 8, 12))
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
