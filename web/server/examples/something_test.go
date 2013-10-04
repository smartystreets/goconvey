package examples

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestPassingStory(t *testing.T) {
	Convey("A passing story", t, func() {
		Convey("This is the coolest thing ever", func() {
			So("This test passes", ShouldContainSubstring, "pass")
			So("This test passes", ShouldContainSubstring, "pass")
		})

		Convey("Hakuna Matata", func() {
			So("1", ShouldEqual, "1")
		})
	})
}

// func TestOldSchool_Panics(t *testing.T) {
// 	if []int{}[0] == 42 {
// 		t.Log("We shouldn't EVER get here.")
// 	}
// }

func TestOldSchool_Passes(t *testing.T) {
	// passes implicitly
}

func TestOldSchool_PassesWithMessage(t *testing.T) {
	t.Log("I am a passing test.\nWith a newline.")
}

// func TestOldSchool_Failure(t *testing.T) {
// 	t.Fail() // no message
// }

// func TestOldSchool_FailureWithReason(t *testing.T) {
// 	t.Error("I am a failing test.")
// }
