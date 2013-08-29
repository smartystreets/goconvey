package goconvey

import (
	"github.com/mdwhatcott/goconvey/convey/execution"
	"testing"
)

func TestNothingInScope(t *testing.T) {
	output := prepare()

	run()

	expect(t, "", output)
}

func TestSingleScope(t *testing.T) {
	output := prepare()

	Convey("hi", func() {
		output += "done"
	})

	run()

	expect(t, "done", output)
}

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	output := prepare()

	Convey("1", func() {
		output += "1"
	})

	Convey("2", func() {
		output += "2"
	})

	run()

	expect(t, "12", output)
}

func TestNestedScopes(t *testing.T) {
	output := prepare()

	Convey("a", func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})
		})
	})

	run()

	expect(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	output := prepare()

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

	run()

	expect(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func TestSingleScopeWithConveyAndNestedReset(t *testing.T) {
	output := prepare()

	Convey("1", func() {
		output += "1"

		Reset(func() {
			output += "a"
		})
	})

	run()

	expect(t, "1a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndReset(t *testing.T) {
	output := prepare()

	Convey("reset after each nested convey", func() {
		Convey("first output", func() {
			output += "1"
		})

		Convey("second output", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})
	})

	run()

	expect(t, "1a2a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndMultipleResets(t *testing.T) {
	output := prepare()

	Convey("each reset is run at end of each nested convey", func() {
		Convey("1", func() {
			output += "1"
		})

		Convey("2", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})

		Reset(func() {
			output += "b"
		})
	})

	run()

	expect(t, "1ab2ab", output)
}

func TestPanicAtHigherLevelScopePreventsChildScopesFromRunning(t *testing.T) {
	output := prepare()

	Convey("This step panics", func() {
		Convey("this should NOT be executed", func() {
			output += "1"
		})

		panic("Hi")
	})

	run()

	expect(t, "", output)
}

func TestPanicInChildScopeDoes_NOT_PreventExecutionOfSiblingScopes(t *testing.T) {
	output := prepare()

	Convey("This is the parent", func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling should execute", func() {
			output += "2"
		})
	})

	run()

	expect(t, "2", output)
}

func TestResetsAreAlwaysExecutedAfterScopePanics(t *testing.T) {
	output := prepare()

	Convey("This is the parent", func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling step does not panic", func() {
			output += "a"

			Reset(func() {
				output += "b"
			})
		})

		Reset(func() {
			output += "2"
		})
	})

	run()

	expect(t, "2ab2", output)
}

func prepare() string {
	execution.SpecRunner = execution.NewScopeRunner()
	return ""
}

func run() {
	execution.SpecRunner.Run()
}

func expect(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
