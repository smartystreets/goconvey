package main

import (
	"encoding/json"
	"github.com/smartystreets/goconvey/reporting"
	"testing"
)

func TestParsePackage_OldSchoolWithFailureOutput(t *testing.T) {
	actual := parsePackageResults(inputOldSchool_Fails)
	assertDeepEqual(t, expectedOldSchool_Fails, *actual)
}

func TestParsePackage_OldSchoolWithSuccessOutput(t *testing.T) {
	actual := parsePackageResults(inputOldSchool_Passes)
	assertDeepEqual(t, expectedOldSchool_Passes, *actual)
}

func assertDeepEqual(t *testing.T, expected, actual interface{}) {
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	if string(a) != string(b) {
		t.Errorf(failureTemplate, string(a), string(b))
	}
}

// TODO: tweak the durations and package names, make each test input unique...

const inputOldSchool_Fails = `
=== RUN TestOldSchool_Passes
--- PASS: TestOldSchool_Passes (0.00 seconds)
=== RUN TestOldSchool_PassesWithMessage
--- PASS: TestOldSchool_PassesWithMessage (0.00 seconds)
	old_school_test.go:10: I am a passing test.
		With a newline.
=== RUN TestOldSchool_Failure
--- FAIL: TestOldSchool_Failure (0.00 seconds)
=== RUN TestOldSchool_FailureWithReason
--- FAIL: TestOldSchool_FailureWithReason (0.00 seconds)
	old_school_test.go:18: I am a failing test.
FAIL
exit status 1
FAIL	github.com/smartystreets/goconvey/webserver/examples	0.017s
`

const inputOldSchool_Passes = `
=== RUN TestOldSchool_Passes
--- PASS: TestOldSchool_Passes (0.00 seconds)
=== RUN TestOldSchool_PassesWithMessage
--- PASS: TestOldSchool_PassesWithMessage (0.00 seconds)
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
			Elapsed:  0.0,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.0,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "I am a passing test.\nWith a newline.",
			Stories:  []reporting.ScopeResult{},
		},
	},
}

var expectedOldSchool_Fails = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/webserver/examples",
	Elapsed:     0.017,
	Passed:      false,
	TestResults: []TestResult{
		TestResult{
			TestName: "TestOldSchool_Passes",
			Elapsed:  0.0,
			Passed:   true,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_PassesWithMessage",
			Elapsed:  0.0,
			Passed:   true,
			File:     "old_school_test.go",
			Line:     10,
			Message:  "I am a passing test.\nWith a newline.",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_Failure",
			Elapsed:  0.0,
			Passed:   false,
			File:     "",
			Line:     0,
			Message:  "",
			Stories:  []reporting.ScopeResult{},
		},
		TestResult{
			TestName: "TestOldSchool_FailureWithReason",
			Elapsed:  0.0,
			Passed:   false,
			File:     "old_school_test.go",
			Line:     18,
			Message:  "I am a failing test.",
			Stories:  []reporting.ScopeResult{},
		},
	},
}

const failureTemplate = "Comparison failed:\n  Expected: %v\n    Actual: %v\n"

/*
func TestParseJsonOutput(t *testing.T) {
	var parsed PackageResult

	Convey("Subject: Parse output from 'go test -json'", t, func() {
		Convey("When the package passed all tests", func() {
			parsed = parsePackageResult(PassingOutput)

			Convey("The parsed result should be correct", func() {
				So(parsed, ShouldResemble, ParsedPassingOutput)
			})
		})
		Convey("When the package failed any tests", func() {
			parsed = parsePackageResult(FailingOutput)

			Convey("The parsed result should be correct", func() {
				So(parsed, ShouldResemble, ParsedFailingOutput)
			})
		})
	})
}

const PassingOutput = `[
  {
    "Title": "TestParseJsonOutput",
    "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/goconvey-server/story_parse_test.go",
    "Line": 11,
    "Depth": 0,
    "Assertions": []
  }
],PASS
ok  	github.com/smartystreets/goconvey/web/goconvey-server	0.031s`

var ParsedPassingOutput = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/web/goconvey-server",
	Elapsed:     .031,
	Passed:      true,
	Stories: []StoryResult{
		[]reporting.ScopeResult{
			reporting.ScopeResult{
				Title:      "TestParseJsonOutput",
				File:       "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/goconvey-server/story_parse_test.go",
				Line:       11,
				Depth:      0,
				Assertions: []reporting.AssertionResult{},
			},
		},
	},
}

const FailingOutput = `[
  {
    "Title": "TestParseJsonOutput",
    "File": "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/goconvey-server/story_parse_test.go",
    "Line": 11,
    "Depth": 0,
    "Assertions": []
  }
],--- FAIL: TestParseJsonOutput (0.00 seconds)
FAIL
exit status 1
FAIL	github.com/smartystreets/goconvey/web/goconvey-server	0.032s`

var ParsedFailingOutput = PackageResult{
	PackageName: "github.com/smartystreets/goconvey/web/goconvey-server",
	Elapsed:     .032,
	Passed:      false,
	Stories: []StoryResult{
		[]reporting.ScopeResult{
			reporting.ScopeResult{
				Title:      "TestParseJsonOutput",
				File:       "/Users/mike/work/dev/goconvey/src/github.com/smartystreets/goconvey/web/goconvey-server/story_parse_test.go",
				Line:       11,
				Depth:      0,
				Assertions: []reporting.AssertionResult{},
			},
		},
	},
}
*/
