package parser

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/reporting"
	"github.com/smartystreets/goconvey/web/goconvey-server/results"
	"strings"
	"testing"
)

func TestParsePackage_NoGoFiles_ReturnsPackageResult(t *testing.T) {
	packageName := expected_NoGoFiles.PackageName
	actual := ParsePackageResults(packageName, input_NoGoFiles)
	assertEqual(t, expected_NoGoFiles, *actual)
}

func TestParsePackage_NoTestFiles_ReturnsPackageResult(t *testing.T) {
	packageName := expected_NoTestFiles.PackageName
	actual := ParsePackageResults(packageName, input_NoTestFiles)
	assertEqual(t, expected_NoTestFiles, *actual)
}

func TestParsePacakge_NoTestFunctions_ReturnsPackageResult(t *testing.T) {
	packageName := expected_NoTestFunctions.PackageName
	actual := ParsePackageResults(packageName, input_NoTestFunctions)
	assertEqual(t, expected_NoTestFunctions, *actual)
}

func TestParsePackage_BuildFailed_ReturnsPackageResult(t *testing.T) {
	packageName := expected_BuildFailed_InvalidPackageDeclaration.PackageName
	actual := ParsePackageResults(packageName, input_BuildFailed_InvalidPackageDeclaration)
	assertEqual(t, expected_BuildFailed_InvalidPackageDeclaration, *actual)

	packageName = expected_BuildFailed_OtherErrors.PackageName
	actual = ParsePackageResults(packageName, input_BuildFailed_OtherErrors)
	assertEqual(t, expected_BuildFailed_OtherErrors, *actual)
}

func TestParsePackage_OldSchoolWithFailureOutput_ReturnsCompletePackageResult(t *testing.T) {
	packageName := expectedOldSchool_Fails.PackageName
	actual := ParsePackageResults(packageName, inputOldSchool_Fails)
	assertEqual(t, expectedOldSchool_Fails, *actual)
}

func TestParsePackage_OldSchoolWithSuccessOutput_ReturnsCompletePackageResult(t *testing.T) {
	packageName := expectedOldSchool_Passes.PackageName
	actual := ParsePackageResults(packageName, inputOldSchool_Passes)
	assertEqual(t, expectedOldSchool_Passes, *actual)
}

func TestParsePackage_OldSchoolWithPanicOutput_ReturnsCompletePackageResult(t *testing.T) {
	packageName := expectedOldSchool_Panics.PackageName
	actual := ParsePackageResults(packageName, inputOldSchool_Panics)
	assertEqual(t, expectedOldSchool_Panics, *actual)
}

func TestParsePackage_GoConveyOutput_ReturnsCompletePackageResult(t *testing.T) {
	packageName := expectedGoConvey.PackageName
	actual := ParsePackageResults(packageName, inputGoConvey)
	assertEqual(t, expectedGoConvey, *actual)
}

func TestParsePackage_ActualPackageNameDifferentThanDirectoryName_ReturnsActualPackageName(t *testing.T) {
	packageName := strings.Replace(expectedGoConvey.PackageName, "examples", "stuff", -1)
	actual := ParsePackageResults(packageName, inputGoConvey)
	assertEqual(t, expectedGoConvey, *actual)
}

func TestParsePackage_GoConveyOutputMalformed_CausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprintf("%v", r)
			if !strings.Contains(message, "bug report") {
				t.Errorf("Should have panicked with a request to file a bug report but we received this error instead: %s", message)
			}
		} else {
			t.Errorf("Should have panicked with a request to file a bug report but we received no error.")
		}
	}()

	ParsePackageResults(expectedGoConvey.PackageName, inputGoConvey_Malformed)
}

func TestParsePackage_GoConveyWithRandomOutput_ReturnsPackageResult(t *testing.T) {
	packageName := expectedGoConvey_WithRandomOutput.PackageName
	actual := ParsePackageResults(packageName, inputGoConvey_WithRandomOutput)
	assertEqual(t, expectedGoConvey_WithRandomOutput, *actual)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	if string(a) != string(b) {
		t.Errorf(failureTemplate, string(a), string(b))
	}
}

const input_NoGoFiles = `can't load package: package github.com/smartystreets/goconvey: no Go source files in /Users/matt/Work/Dev/goconvey/src/github.com/smartystreets/goconvey`

var expected_NoGoFiles = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey",
	Outcome:     results.NoGoFiles,
	BuildOutput: input_NoGoFiles,
	TestResults: []results.TestResult{},
}

const input_NoTestFiles = `?   	pkg.smartystreets.net/liveaddress-zipapi	[no test files]`

var expected_NoTestFiles = results.PackageResult{
	PackageName: "pkg.smartystreets.net/liveaddress-zipapi",
	Outcome:     results.NoTestFiles,
	BuildOutput: input_NoTestFiles,
	TestResults: []results.TestResult{},
}

const input_NoTestFunctions = `testing: warning: no tests to run`

var expected_NoTestFunctions = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Outcome:     results.NoTestFunctions,
	BuildOutput: input_NoTestFunctions,
	TestResults: []results.TestResult{},
}

const input_BuildFailed_InvalidPackageDeclaration = `
can't load package: package github.com/smartystreets/goconvey/examples:
bowling_game_test.go:9:1: expected 'package', found 'IDENT' asdf
bowling_game_test.go:10:1: invalid package name _
`

var expected_BuildFailed_InvalidPackageDeclaration = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/examples",
	Outcome:     results.BuildFailure,
	BuildOutput: strings.TrimSpace(input_BuildFailed_InvalidPackageDeclaration),
	TestResults: []results.TestResult{},
}

const input_BuildFailed_OtherErrors = `
# github.com/smartystreets/goconvey/examples
./bowling_game_test.go:22: undefined: game
./bowling_game_test.go:22: cannot assign to game
./bowling_game_test.go:25: undefined: game
./bowling_game_test.go:28: undefined: game
./bowling_game_test.go:33: undefined: game
./bowling_game_test.go:36: undefined: game
./bowling_game_test.go:41: undefined: game
./bowling_game_test.go:42: undefined: game
./bowling_game_test.go:43: undefined: game
./bowling_game_test.go:46: undefined: game
./bowling_game_test.go:46: too many errors
FAIL	github.com/smartystreets/goconvey/examples [build failed]
`

var expected_BuildFailed_OtherErrors = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/examples",
	Outcome:     results.BuildFailure,
	BuildOutput: strings.TrimSpace(input_BuildFailed_OtherErrors),
	TestResults: []results.TestResult{},
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

var expectedOldSchool_Passes = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.018,
	Outcome:     results.Passed,
	TestResults: []results.TestResult{
		results.TestResult{
			TestName: "TestOldSchool_Passes",
			Elapsed:  0.02,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		results.TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.05,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "old_school_test.go:10: I am a passing test.\nWith a newline.",
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

var expectedOldSchool_Fails = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Outcome:     results.Failed,
	Elapsed:     0.017,
	TestResults: []results.TestResult{
		results.TestResult{
			TestName: "TestOldSchool_Passes",
			Elapsed:  0.01,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		results.TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.03,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "old_school_test.go:10: I am a passing test.\nWith a newline.",
			Stories:  []reporting.ScopeResult{},
		},
		results.TestResult{
			TestName: "TestOldSchool_Failure",
			Elapsed:  0.06,
			Passed:   false,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		results.TestResult{
			TestName: "TestOldSchool_FailureWithReason",
			Elapsed:  0.11,
			Passed:   false,
			File:     "old_school_test.go",
			Line:     18,
			Message:  "old_school_test.go:18: I am a failing test.",
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

var expectedOldSchool_Panics = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.014,
	Outcome:     results.Panicked,
	TestResults: []results.TestResult{
		results.TestResult{
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

const inputGoConvey_Malformed = `
=== RUN TestPassingStory
>>>>>
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

      ;aiwheopinen39 n3902n92m

      "Error": null,
      "Skipped": false,
      "StackTrace": "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/webserver/examples.func·001()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:10 +0xe3\ngithub.com/smartystreets/goconvey/webserver/examples.TestPassingStory(0x210314000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/webserver/examples/old_school_test.go:11 +0xec\ntesting.tRunner(0x210314000, 0x21ab10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n"
    }
  ]
},
<<<<<
--- PASS: TestPassingStory (0.01 seconds)
PASS
ok  	github.com/smartystreets/goconvey/webserver/examples	0.019s
`

const inputGoConvey = `
=== RUN TestPassingStory
>>>>>
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
<<<<<
--- PASS: TestPassingStory (0.01 seconds)
PASS
ok  	github.com/smartystreets/goconvey/webserver/examples	0.019s
`

var expectedGoConvey = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.019,
	Outcome:     results.Passed,
	TestResults: []results.TestResult{
		results.TestResult{
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

const inputGoConvey_WithRandomOutput = `
=== RUN TestPassingStory
*** Hello, World! (1) ***
*** Hello, World! (2) ***
*** Hello, World! (3) ***
>>>>>
{
  "Title": "A passing story",
  "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
  "Line": 16,
  "Depth": 0,
  "Assertions": [
    {
      "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
      "Line": 14,
      "Failure": "",
      "Error": null,
      "Skipped": false,
      "StackTrace": "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/web/server/testing.func·001()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:14 +0x186\ngithub.com/smartystreets/goconvey/web/server/testing.TestPassingStory(0x210315000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:16 +0x1b9\ntesting.tRunner(0x210315000, 0x21bb10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n"
    }
  ]
},
<<<<<
*** Hello, World! (4)***
*** Hello, World! (5) ***
>>>>>
{
  "Title": "A passing story",
  "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
  "Line": 22,
  "Depth": 0,
  "Assertions": [
    {
      "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
      "Line": 20,
      "Failure": "",
      "Error": null,
      "Skipped": false,
      "StackTrace": "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/web/server/testing.func·002()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:20 +0x186\ngithub.com/smartystreets/goconvey/web/server/testing.TestPassingStory(0x210315000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:22 +0x294\ntesting.tRunner(0x210315000, 0x21bb10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n"
    }
  ]
},
<<<<<
*** Hello, World! (6) ***
--- PASS: TestPassingStory (0.03 seconds)
PASS
ok  	github.com/smartystreets/goconvey/web/server/testing	0.024s
`

var expectedGoConvey_WithRandomOutput = results.PackageResult{
	PackageName: "github.com/smartystreets/goconvey/web/server/testing",
	Elapsed:     0.024,
	Outcome:     results.Passed,
	TestResults: []results.TestResult{
		results.TestResult{
			TestName: "TestPassingStory",
			Elapsed:  0.03,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "*** Hello, World! (1) ***\n*** Hello, World! (2) ***\n*** Hello, World! (3) ***\n*** Hello, World! (4)***\n*** Hello, World! (5) ***\n*** Hello, World! (6) ***",
			Stories: []reporting.ScopeResult{
				reporting.ScopeResult{
					Title: "A passing story",
					File:  "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
					Line:  16,
					Depth: 0,
					Assertions: []reporting.AssertionResult{
						reporting.AssertionResult{
							File:       "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
							Line:       14,
							Failure:    "",
							Error:      nil,
							Skipped:    false,
							StackTrace: "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/web/server/testing.func·001()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:14 +0x186\ngithub.com/smartystreets/goconvey/web/server/testing.TestPassingStory(0x210315000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:16 +0x1b9\ntesting.tRunner(0x210315000, 0x21bb10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n",
						},
					},
				},
				reporting.ScopeResult{
					Title: "A passing story",
					File:  "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
					Line:  22,
					Depth: 0,
					Assertions: []reporting.AssertionResult{
						reporting.AssertionResult{
							File:       "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go",
							Line:       20,
							Failure:    "",
							Error:      nil,
							Skipped:    false,
							StackTrace: "goroutine 3 [running]:\ngithub.com/smartystreets/goconvey/web/server/testing.func·002()\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:20 +0x186\ngithub.com/smartystreets/goconvey/web/server/testing.TestPassingStory(0x210315000)\n\u0009/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/server/testing/go_test.go:22 +0x294\ntesting.tRunner(0x210315000, 0x21bb10)\n\u0009/usr/local/go/src/pkg/testing/testing.go:353 +0x8a\ncreated by testing.RunTests\n\u0009/usr/local/go/src/pkg/testing/testing.go:433 +0x86b\n",
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

GoConvey style tests with random output:

	package examples

	import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
	)

	func TestPassingStory(t *testing.T) {
		fmt.Println("*** Hello, World! (1) ***")

		Convey("A passing story", t, func() {
			fmt.Println("*** Hello, World! (2) ***")
			So("This test passes", ShouldContainSubstring, "pass")
			fmt.Println("*** Hello, World! (3) ***")
		})

		Convey("A passing story", t, func() {
			fmt.Println("*** Hello, World! (4)***")
			So("This test passes", ShouldContainSubstring, "pass")
			fmt.Println("*** Hello, World! (5) ***")
		})

		fmt.Println("*** Hello, World! (6) ***")
	}


*/
