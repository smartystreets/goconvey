package examples

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPassingStory(t *testing.T) {
	Convey("A passing story", t, func() {
		So("This test passes", ShouldContainSubstring, "pass")
	})
}

func TestOldSchool_Panics(t *testing.T) {
	if []int{}[0] == 42 {
		t.Log("We shouldn't EVER get here.")
	}
}

func TestOldSchool_Passes(t *testing.T) {
	// passes implicitly
}

func TestOldSchool_PassesWithMessage(t *testing.T) {
	t.Log("I am a passing test.\nWith a newline.")
}

func TestOldSchool_Failure(t *testing.T) {
	t.Fail() // no message
}

func TestOldSchool_FailureWithReason(t *testing.T) {
	t.Error("I am a failing test.")
}
