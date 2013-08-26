package goconvey_test

import (
	"github.com/mdwhatcott/goconvey"
	"testing"
)

func TestNothingInScope(t *testing.T) {
	runner := goconvey.NewRunner()
	output := ""

	runner.Run()

	expect(t, "", output)
}

func TestSingleScopeWithConvey(t *testing.T) {
	runner := goconvey.NewRunner()
	output := ""

	runner.Register("hi", func() {
		output += "done"
	})

	runner.Run()

	expect(t, "done", output)
}

// func TestSingleScopeWithConveyAndReset(t *testing.T) {
// 	runner := goconvey.NewRunner()
// 	output := ""

// 	runner.Register("1", func() {
// 		output += "1"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.Run()

// 	expect(t, "1a", output)
// }

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	runner := goconvey.NewRunner()
	output := ""

	runner.Register("1", func() {
		output += "1"
	})

	runner.Register("2", func() {
		output += "2"
	})

	runner.Run()

	expect(t, "12", output)
}

// func TestSingleScopeWithMultipleConveysAndReset(t *testing.T) {
// 	runner := goconvey.NewRunner()
// 	output := ""

// 	runner.Register("1", func() {
// 		output += "1"
// 	})

// 	runner.Register("2 again", func() {
// 		output += "2"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.Run()

// 	expect(t, "12a", output)
// }

// func TestSingleScopeWithMultipleConveysAndMultipleResets(t *testing.T) {
// 	runner := goconvey.NewRunner()
// 	output := ""

// 	runner.Register("1", func() {
// 		output += "1"
// 	})

// 	runner.Register("2", func() {
// 		output += "2"
// 	})

// 	runner.Reset("a", func() {
// 		output += "a"
// 	})

// 	runner.Reset("b", func() {
// 		output += "b"
// 	})

// 	runner.Run()

// 	expect(t, "12ab", output)
// }

func TestNestedScopes(t *testing.T) {
	runner := goconvey.NewRunner()
	output := ""

	runner.Register("a", func() {
		output += "a "

		runner.Register("aa", func() {
			output += "aa "

			runner.Register("aaa", func() {
				output += "aaa | "
			})
		})
	})

	runner.Run()

	expect(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	runner := goconvey.NewRunner()
	output := ""

	runner.Register("a", func() {
		output += "a "

		runner.Register("aa", func() {
			output += "aa "

			runner.Register("aaa", func() {
				output += "aaa | "
			})

			runner.Register("aaa1", func() {
				output += "aaa1 | "
			})
		})

		runner.Register("ab", func() {
			output += "ab "

			runner.Register("abb", func() {
				output += "abb | "
			})
		})
	})

	runner.Run()

	expect(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func expect(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
