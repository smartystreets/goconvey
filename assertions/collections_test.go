package assertions

import (
	"testing"
)

func TestShouldContain(t *testing.T) {
	fail(t, so([]int{}, ShouldContain), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so([]int{}, ShouldContain, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(Thing1{}, ShouldContain, 1), "You must provide a valid collection (was assertions.Thing1)!")
	fail(t, so(nil, ShouldContain, 1), "You must provide a valid collection (was <nil>)!")
	fail(t, so([]int{1}, ShouldContain, 2), "Expected the collection ([]int) to contain member: '2' (but it didn't)!")

	pass(t, so([]int{1}, ShouldContain, 1))
	pass(t, so([]int{1, 2, 3}, ShouldContain, 2))
}

func TestShouldNotContain(t *testing.T) {
	fail(t, so([]int{}, ShouldNotContain), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so([]int{}, ShouldNotContain, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(Thing1{}, ShouldNotContain, 1), "You must provide a valid collection (was assertions.Thing1)!")
	fail(t, so(nil, ShouldNotContain, 1), "You must provide a valid collection (was <nil>)!")

	fail(t, so([]int{1}, ShouldNotContain, 1), "Expected the collection ([]int) NOT to contain member: '1' (but it did)!")
	fail(t, so([]int{1, 2, 3}, ShouldNotContain, 2), "Expected the collection ([]int) NOT to contain member: '2' (but it did)!")

	pass(t, so([]int{1}, ShouldNotContain, 2))
}
