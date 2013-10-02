package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/reporting"
	"testing"
)

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
