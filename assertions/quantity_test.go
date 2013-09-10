package assertions

import (
	"runtime"
	"testing"
)

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
