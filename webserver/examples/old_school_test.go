package examples

import "testing"

func TestOldSchool_Passes(t *testing.T) {
	// passes implicitly
}

func TestOldSchool_PassesWithMessage(t *testing.T) {
	t.Log("I am a passing test.\nWith a newline.")
}

/*
func TestOldSchool_Failure(t *testing.T) {
	t.Fail() // no message
}

func TestOldSchool_FailureWithReason(t *testing.T) {
	t.Error("I am a failing test.")
}
*/
