package assertions

import "testing"

func TestShouldHaveSameTypeAs(t *testing.T) {
	serializer = newFakeSerializer()

	fail(t, so(1, ShouldHaveSameTypeAs), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(1, ShouldHaveSameTypeAs, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(nil, ShouldHaveSameTypeAs, 0), "int|<nil>|Expected '<nil>' to be: 'int' (but was: '<nil>')!")
	fail(t, so(1, ShouldHaveSameTypeAs, "asdf"), "string|int|Expected '1' to be: 'string' (but was: 'int')!")

	pass(t, so(1, ShouldHaveSameTypeAs, 0))
	pass(t, so(nil, ShouldHaveSameTypeAs, nil))
}

func TestShouldNotHaveSameTypeAs(t *testing.T) {
	fail(t, so(1, ShouldNotHaveSameTypeAs), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(1, ShouldNotHaveSameTypeAs, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(1, ShouldNotHaveSameTypeAs, 0), "Expected '1' to NOT be: 'int' (but it was)!")
	fail(t, so(nil, ShouldNotHaveSameTypeAs, nil), "Expected '<nil>' to NOT be: '<nil>' (but it was)!")

	pass(t, so(nil, ShouldNotHaveSameTypeAs, 0))
	pass(t, so(1, ShouldNotHaveSameTypeAs, "asdf"))
}
