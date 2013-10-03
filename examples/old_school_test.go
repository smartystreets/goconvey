package examples

import "testing"

func TestOldSchool_Passes(t *testing.T) {
	// pass
}

func TestOldSchool_Failure(t *testing.T) {
	t.Fail()
}

func TestOldSchool_FailureWithReason(t *testing.T) {
	t.Error("I am a failing test.")
}
