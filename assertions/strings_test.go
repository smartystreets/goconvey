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
	pass(t, so("superman", ShouldNotStartWith, "man"))

	fail(t, so(1, ShouldNotStartWith, 2), "Both arguments to this assertions must be strings (you provided int and int).")
}

func TestShouldEndWith(t *testing.T) {
	fail(t, so("", ShouldEndWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so("", ShouldEndWith, "", ""), "This assertion requires exactly 1 comparison values (you provided 2).")

	pass(t, so("", ShouldEndWith, ""))
	pass(t, so("superman", ShouldEndWith, "man"))
	fail(t, so("superman", ShouldEndWith, "super"), "Expected 'superman' to end with 'super' (but it didn't)!")
	fail(t, so("superman", ShouldEndWith, "blah"), "Expected 'superman' to end with 'blah' (but it didn't)!")

	fail(t, so(1, ShouldEndWith, 2), "Both arguments to this assertions must be strings (you provided int and int).")
}

func TestShouldNotEndWith(t *testing.T) {
	fail(t, so("", ShouldNotEndWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so("", ShouldNotEndWith, "", ""), "This assertion requires exactly 1 comparison values (you provided 2).")

	fail(t, so("", ShouldNotEndWith, ""), "Expected '<empty>' NOT to end with '<empty>' (but it did)!")
	fail(t, so("superman", ShouldNotEndWith, "man"), "Expected 'superman' NOT to end with 'man' (but it did)!")
	pass(t, so("superman", ShouldNotEndWith, "super"))

	fail(t, so(1, ShouldNotEndWith, 2), "Both arguments to this assertions must be strings (you provided int and int).")
}
