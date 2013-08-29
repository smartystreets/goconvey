package execution

import (
	"testing"
)

func TestNothingInScope(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Run()

	expect(t, "", output)
}

func TestSingleScope(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("hi", func() {
		output += "done"
	})

	specRunner.Run()

	expect(t, "done", output)
}

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("1", func() {
		output += "1"
	})

	specRunner.Register("2", func() {
		output += "2"
	})

	specRunner.Run()

	expect(t, "12", output)
}

func TestNestedScopes(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("a", func() {
		output += "a "

		specRunner.Register("aa", func() {
			output += "aa "

			specRunner.Register("aaa", func() {
				output += "aaa | "
			})
		})
	})

	specRunner.Run()

	expect(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("a", func() {
		output += "a "

		specRunner.Register("aa", func() {
			output += "aa "

			specRunner.Register("aaa", func() {
				output += "aaa | "
			})

			specRunner.Register("aaa1", func() {
				output += "aaa1 | "
			})
		})

		specRunner.Register("ab", func() {
			output += "ab "

			specRunner.Register("abb", func() {
				output += "abb | "
			})
		})
	})

	specRunner.Run()

	expect(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func TestSingleScopeWithConveyAndNestedReset(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("1", func() {
		output += "1"

		specRunner.RegisterReset(func() {
			output += "a"
		})
	})

	specRunner.Run()

	expect(t, "1a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndReset(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("reset after each nested convey", func() {
		specRunner.Register("first output", func() {
			output += "1"
		})

		specRunner.Register("second output", func() {
			output += "2"
		})

		specRunner.RegisterReset(func() {
			output += "a"
		})
	})

	specRunner.Run()

	expect(t, "1a2a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndMultipleResets(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("each reset is run at end of each nested convey", func() {
		specRunner.Register("1", func() {
			output += "1"
		})

		specRunner.Register("2", func() {
			output += "2"
		})

		specRunner.RegisterReset(func() {
			output += "a"
		})

		specRunner.RegisterReset(func() {
			output += "b"
		})
	})

	specRunner.Run()

	expect(t, "1ab2ab", output)
}

func TestPanicAtHigherLevelScopePreventsChildScopesFromRunning(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("This step panics", func() {
		specRunner.Register("this should NOT be executed", func() {
			output += "1"
		})

		panic("Hi")
	})

	specRunner.Run()

	expect(t, "", output)
}

func TestPanicInChildScopeDoes_NOT_PreventExecutionOfSiblingScopes(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("This is the parent", func() {
		specRunner.Register("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		specRunner.Register("This sibling should execute", func() {
			output += "2"
		})
	})

	specRunner.Run()

	expect(t, "2", output)
}

func TestResetsAreAlwaysExecutedAfterScopePanics(t *testing.T) {
	specRunner, output := prepare()

	specRunner.Register("This is the parent", func() {
		specRunner.Register("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		specRunner.Register("This sibling step does not panic", func() {
			output += "a"

			specRunner.RegisterReset(func() {
				output += "b"
			})
		})

		specRunner.RegisterReset(func() {
			output += "2"
		})
	})

	specRunner.Run()

	expect(t, "2ab2", output)
}

func prepare() (*runner, string) {
	return NewSpecRunner(), ""
}

func expect(t *testing.T, expected, actual string) {
	if actual != expected {
		t.Errorf("Expected '%s', got '%s'", expected, actual)
	}
}
