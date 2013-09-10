package assertions

import (
	"testing"
)

func TestShouldStartWith(t *testing.T) {
	fail(t, so("", ShouldStartWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so("", ShouldStartWith, "asdf", "asdf"), "This assertion requires exactly 1 comparison values (you provided 2).")

	pass(t, so("", ShouldStartWith, ""))
	pass(t, so("superman", ShouldStartWith, "super"))
	fail(t, so("superman", ShouldStartWith, "bat"), "Expected 'superman' to start with 'bat' (but it didn't)!")
	fail(t, so("superman", ShouldStartWith, "man"), "Expected 'superman' to start with 'man' (but it didn't)!")

	fail(t, so(1, ShouldStartWith, 2), "Both arguments to this assertions must be strings (you provided int and int).")
}

func TestShouldNotStartWith(t *testing.T) {
	fail(t, so("", ShouldNotStartWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so("", ShouldNotStartWith, "asdf", "asdf"), "This assertion requires exactly 1 comparison values (you provided 2).")

	fail(t, so("", ShouldNotStartWith, ""), "Expected '<empty>' NOT to start with '<empty>' (but it did)!")
	fail(t, so("superman", ShouldNotStartWith, "super"), "Expected 'superman' NOT to start with 'super' (but it did)!")
	pass(t, so("superman", ShouldNotStartWith, "bat"))

	fail(t, so(1, ShouldNotStartWith, 2), "Both arguments to this assertions must be strings (you provided int and int).")
}
