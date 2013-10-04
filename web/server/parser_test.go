package main

import (
	"encoding/json"
	"github.com/smartystreets/goconvey/reporting"
	"strings"
	"testing"
)

func TestParsePackage_OldSchoolWithFailureOutput(t *testing.T) {
	actual := parsePackageResults(inputOldSchool_Fails)
	assertEqual(t, expectedOldSchool_Fails, *actual)
}

func TestParsePackage_OldSchoolWithSuccessOutput(t *testing.T) {
	actual := parsePackageResults(inputOldSchool_Passes)
	assertEqual(t, expectedOldSchool_Passes, *actual)
}

func TestParsePackage_OldSchoolWithPanicOutput(t *testing.T) {
	actual := parsePackageResults(inputOldSchool_Panics)
	assertEqual(t, expectedOldSchool_Panics, *actual)
}

func TestParsePackage_GoConveyOutput(t *testing.T) {
	actual := parsePackageResults(inputGoConvey)
	assertEqual(t, expectedGoConvey, *actual)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	if string(a) != string(b) {
		t.Errorf(failureTemplate, string(a), string(b))
	}
}

const inputOldSchool_Passes = `
=== RUN TestOldSchool_Passes
--- PASS: TestOldSchool_Passes (0.02 seconds)
=== RUN TestOldSchool_PassesWithMessage
--- PASS: TestOldSchool_PassesWithMessage (0.05 seconds)
	old_school_test.go:10: I am a passing test.
		With a newline.
PASS
ok  	github.com/smartystreets/goconvey/webserver/examples	0.018s
`

var expectedOldSchool_Passes = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.018,
	Passed:      true,
	TestResults: []TestResult{
		TestResult{
			TestName: "TestOldSchool_Passes",
			Elapsed:  0.02,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.05,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "I am a passing test.\nWith a newline.",
			Stories:  []reporting.ScopeResult{},
		},
	},
}

const inputOldSchool_Fails = `
=== RUN TestOldSchool_Passes
--- PASS: TestOldSchool_Passes (0.01 seconds)
=== RUN TestOldSchool_PassesWithMessage
--- PASS: TestOldSchool_PassesWithMessage (0.03 seconds)
	old_school_test.go:10: I am a passing test.
		With a newline.
=== RUN TestOldSchool_Failure
--- FAIL: TestOldSchool_Failure (0.06 seconds)
=== RUN TestOldSchool_FailureWithReason
--- FAIL: TestOldSchool_FailureWithReason (0.11 seconds)
	old_school_test.go:18: I am a failing test.
FAIL
exit status 1
FAIL	github.com/smartystreets/goconvey/webserver/examples	0.017s
`

var expectedOldSchool_Fails = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.017,
	Passed:      false,
	TestResults: []TestResult{
		TestResult{
			TestName: "TestOldSchool_Passes",
			Elapsed:  0.01,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.03,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "I am a passing test.\nWith a newline.",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_Failure",
			Elapsed:  0.06,
			Passed:   false,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_FailureWithReason",
			Elapsed:  0.11,
			Passed:   false,
			File:     "old_school_test.go",
			Line:     18,
			Message:  "I am a failing test.",
			Stories:  []reporting.ScopeResult{},
		},
	},
}

const inputOldSchool_Panics = `
=== RUN TestOldSchool_Panics
--- FAIL: TestOldSchool_Panics (0.02 seconds)
panic: runtime error: index out of range [recovered]
	panic: runtime error: index out of range

goroutine 3 [running]:
testing.func·004()
	/usr/local/go/src/pkg/testing/testing.go:348 +0xcd
github.com/smartystreets/goconvey/webserver/examples.TestOldSchool_Panics(0x210292000)
	/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/something_test.go:15 +0xec
testing.tRunner(0x210292000, 0x1b09f0)
	/usr/local/go/src/pkg/testing/testing.go:353 +0x8a
created by testing.RunTests
	/usr/local/go/src/pkg/testing/testing.go:433 +0x86b

goroutine 1 [chan receive]:
testing.RunTests(0x138f38, 0x1b09f0, 0x1, 0x1, 0x1, ...)
	/usr/local/go/src/pkg/testing/testing.go:434 +0x88e
testing.Main(0x138f38, 0x1b09f0, 0x1, 0x1, 0x1b7f60, ...)
	/usr/local/go/src/pkg/testing/testing.go:365 +0x8a
main.main()
	github.com/smartystreets/goconvey/webserver/examples/_test/_testmain.go:43 +0x9a
exit status 2
FAIL	github.com/smartystreets/goconvey/webserver/examples	0.014s
`

var expectedOldSchool_Panics = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.014,
	Passed:      false,
	TestResults: []TestResult{
		TestResult{
			TestName: "TestOldSchool_Panics",
			Elapsed:  0.02,
			Passed:   false,
			File:     "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/something_test.go",
			Line:     15,
			Message:  "",
			Error: strings.Replace(`panic: runtime error: index out of range [recovered]
	panic: runtime error: index out of range

goroutine 3 [running]:
testing.func·004()
	/usr/local/go/src/pkg/testing/testing.go:348 +0xcd
github.com/smartystreets/goconvey/webserver/examples.TestOldSchool_Panics(0x210292000)
	/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/something_test.go:15 +0xec
testing.tRunner(0x210292000, 0x1b09f0)
	/usr/local/go/src/pkg/testing/testing.go:353 +0x8a
created by testing.RunTests
	/usr/local/go/src/pkg/testing/testing.go:433 +0x86b

goroutine 1 [chan receive]:
testing.RunTests(0x138f38, 0x1b09f0, 0x1, 0x1, 0x1, ...)
	/usr/local/go/src/pkg/testing/testing.go:434 +0x88e
testing.Main(0x138f38, 0x1b09f0, 0x1, 0x1, 0x1b7f60, ...)
	/usr/local/go/src/pkg/testing/testing.go:365 +0x8a
main.main()
	github.com/smartystreets/goconvey/webserver/examples/_test/_testmain.go:43 +0x9a`, "\u0009", "\t", -1),
			Stories: []reporting.ScopeResult{},
		},
	},
}

const inputGoConvey = `
=== RUN TestPassingStory
{
  "Title": "A passing story",
  "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go",
  "Line": 11,
  "Depth": 0,
  "Assertions": [
    {
      "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go",
      "Line": 10,
      "Failure": "",
      "Error": null,
      "Skipped": false,
      "StackTrace": "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/webserver/examples.func·001()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:10 +0xe3\ngithub.com/smartystreets/goconvey/webserver/examples.TestPassingStory(0x210314000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:11 +0xec\ntesting.tRunner(0x210314000, 0x21ab10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n"
    }
  ]
},
--- PASS: TestPassingStory (0.01 seconds)
PASS
ok  	github.com/smartystreets/goconvey/webserver/examples	0.019s
`

var expectedGoConvey = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.019,
	Passed:      true,
	TestResults: []TestResult{
		TestResult{
			TestName: "TestPassingStory",
			Elapsed:  0.01,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories: []reporting.ScopeResult{
				reporting.ScopeResult{
					Title: "A passing story",
					File:  "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go",
					Line:  11,
					Depth: 0,
					Assertions: []reporting.AssertionResult{
						reporting.AssertionResult{
							File:       "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go",
							Line:       10,
							Failure:    "",
							Error:      nil,
							Skipped:    false,
							StackTrace: "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/webserver/examples.func·001()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:10 +0xe3\ngithub.com/smartystreets/goconvey/webserver/examples.TestPassingStory(0x210314000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:11 +0xec\ntesting.tRunner(0x210314000, 0x21ab10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n",
						},
					},
				},
			},
		},
	},
}

const failureTemplate = "Comparison failed:\n  Expected: %v\n    Actual: %v\n"

/*
Test output for these tests was generated from the following test code:

Old School style tests:

	package examples

	import "testing"

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

GoConvey style tests:

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

*/
