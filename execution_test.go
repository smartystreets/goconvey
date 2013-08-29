package goconvey

import (
	"testing"
)

func TestNothingInScope(t *testing.T) {
	specRunner = newSpecRunner()
	output := ""

	specRunner.run()

	expect(t, "", output)
}

func TestSingleScopeWithConvey(t *testing.T) {
	specRunner = newSpecRunner()
	output := ""

	Convey("hi", t, func() {
		output += "done"
	})

	specRunner.run()

	expect(t, "done", output)
}

// func TestSingleScopeWithConveyAndNestedReset(t *testing.T) {
// 	specRunner = newSpecRunner()
// 	output := ""

// 	Convey("1", t, func() {
// 		output += "1"

// 		Reset("a", func() {
// 			output += "a"
// 		})
// 	})

// 	specRunner.run()

// 	expect(t, "1a", output)
// }

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	specRunner = newSpecRunner()
	output := ""

	Convey("1", t, func() {
		output += "1"
	})

	Convey("2", t, func() {
		output += "2"
	})

	specRunner.run()

	expect(t, "12", output)
}

// func TestSingleScopeWithMultipleConveysAndReset(t *testing.T) {
// 	specRunner = newSpecRunner()
// 	output := ""

// 	Convey("reset after each nested convey", t, func() {
// 		Convey("first output", func() {
// 			output += "1"
// 		})

// 		Convey("second output", t, func() {
// 			output += "2"
// 		})

// 		Reset("a", func() {
// 			output += "a"
// 		})
// 	})

// 	specRunner.run()

// 	expect(t, "1a2a", output)
// }

// func TestSingleScopeWithMultipleConveysAndMultipleResets(t *testing.T) {
// 	specRunner = newSpecRunner()
// 	output := ""

// 	Convey("each reset is run at end of each nested convey", t, func() {
// 		Convey("1", func() {
// 			output += "1"
// 		})

// 		Convey("2", func() {
// 			output += "2"
// 		})

// 		Reset("a", func() {
// 			output += "a"
// 		})

// 		Reset("b", func() {
// 			output += "b"
// 		})
// 	})

// 	specRunner.run()

// 	expect(t, "1ab2ab", output)
// }

func TestNestedScopes(t *testing.T) {
	specRunner = newSpecRunner()
	output := ""

	Convey("a", func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})
		})
	})

	specRunner.run()

	expect(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	specRunner = newSpecRunner()
	output := ""

	Convey("a", func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})

			Convey("aaa1", func() {
				output += "aaa1 | "
			})
		})

		Convey("ab", func() {
			output += "ab "

			Convey("abb", func() {
				output += "abb | "
			})
		})
	})

	specRunner.run()

	expect(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func expect(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
