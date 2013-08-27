package goconvey

import (
	"testing"
)

type FakeTest struct{}

func (self FakeTest) Fail() {}

var test = FakeTest{}

func TestNothingInScope(t *testing.T) {
	runner := newSpecRunner()
	output := ""

	runner.run()

	expect(t, "", output)
}

func TestSingleScopeWithConvey(t *testing.T) {
	runner := newSpecRunner()
	output := ""

	runner.register("hi", func() {
		output += "done"
	})

	runner.run()

	expect(t, "done", output)
}

// func TestSingleScopeWithConveyAndReset(t *testing.T) {
// 	runner := newSpecRunner()
// 	output := ""

// 	runner.register("1", func() {
// 		output += "1"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.run()

// 	expect(t, "1a", output)
// }

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	runner := newSpecRunner()
	output := ""

	runner.register("1", func() {
		output += "1"
	})

	runner.register("2", func() {
		output += "2"
	})

	runner.run()

	expect(t, "12", output)
}

// func TestSingleScopeWithMultipleConveysAndReset(t *testing.T) {
// 	runner := newSpecRunner()
// 	output := ""

// 	runner.register("1", func() {
// 		output += "1"
// 	})

// 	runner.register("2 again", func() {
// 		output += "2"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.run()

// 	expect(t, "12a", output)
// }

// func TestSingleScopeWithMultipleConveysAndMultipleResets(t *testing.T) {
// 	runner := newSpecRunner()
// 	output := ""

// 	runner.register("1", func() {
// 		output += "1"
// 	})

// 	runner.register("2", func() {
// 		output += "2"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.Reset("b", func() {
// 		output += "b"
// 	})

// 	runner.run()

// 	expect(t, "12ab", output)
// }

func TestNestedScopes(t *testing.T) {
	runner := newSpecRunner()
	output := ""

	runner.register("a", func() {
		output += "a "

		runner.register("aa", func() {
			output += "aa "

			runner.register("aaa", func() {
				output += "aaa | "
			})
		})
	})

	runner.run()

	expect(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	runner := newSpecRunner()
	output := ""

	runner.register("a", func() {
		output += "a "

		runner.register("aa", func() {
			output += "aa "

			runner.register("aaa", func() {
				output += "aaa | "
			})

			runner.register("aaa1", func() {
				output += "aaa1 | "
			})
		})

		runner.register("ab", func() {
			output += "ab "

			runner.register("abb", func() {
				output += "abb | "
			})
		})
	})

	runner.run()

	expect(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func expect(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
