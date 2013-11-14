package convey

import (
	"fmt"
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestSingleScopeReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		So(1, ShouldEqual, 1)
	})

	expectEqual(t, "Begin|A|Success|Exit|End", myReporter.wholeStory())
}

func TestNestedScopeReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		Convey("B", func() {
			So(1, ShouldEqual, 1)
		})
	})

	expectEqual(t, "Begin|A|B|Success|Exit|Exit|End", myReporter.wholeStory())
}

func TestFailureReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		So(1, ShouldBeNil)
	})

	expectEqual(t, "Begin|A|Failure|Exit|End", myReporter.wholeStory())
}

func TestComparisonFailureDeserializedAndReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		So("hi", ShouldEqual, "bye")
	})

	expectEqual(t, "Begin|A|Failure(bye/hi)|Exit|End", myReporter.wholeStory())
}

func TestNestedFailureReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		Convey("B", func() {
			So(2, ShouldBeNil)
		})
	})

	expectEqual(t, "Begin|A|B|Failure|Exit|Exit|End", myReporter.wholeStory())
}

func TestSuccessAndFailureReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		So(1, ShouldBeNil)
		So(nil, ShouldBeNil)
	})

	expectEqual(t, "Begin|A|Failure|Success|Exit|End", myReporter.wholeStory())
}

func TestIncompleteActionReportedAsSkipped(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		Convey("B", nil)
	})

	expectEqual(t, "Begin|A|B|Skipped|Exit|Exit|End", myReporter.wholeStory())
}

func TestSkippedConveyReportedAsSkipped(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		SkipConvey("B", func() {
			So(1, ShouldEqual, 1)
		})
	})

	expectEqual(t, "Begin|A|B|Skipped|Exit|Exit|End", myReporter.wholeStory())
}

func TestMultipleSkipsAreReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		Convey("0", func() {
			So(nil, ShouldBeNil)
		})

		SkipConvey("1", func() {})
		SkipConvey("2", func() {})

		Convey("3", nil)
		Convey("4", nil)

		Convey("5", func() {
			So(nil, ShouldBeNil)
		})
	})

	expected := "Begin" +
		"|A|0|Success|Exit|Exit" +
		"|A|1|Skipped|Exit|Exit" +
		"|A|2|Skipped|Exit|Exit" +
		"|A|3|Skipped|Exit|Exit" +
		"|A|4|Skipped|Exit|Exit" +
		"|A|5|Success|Exit|Exit" +
		"|End"

	expectEqual(t, expected, myReporter.wholeStory())
}

func TestSkippedAssertionIsNotReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		SkipSo(1, ShouldEqual, 1)
	})

	expectEqual(t, "Begin|A|Skipped|Exit|End", myReporter.wholeStory())
}

func TestMultipleSkippedAssertionsAreNotReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		SkipSo(1, ShouldEqual, 1)
		So(1, ShouldEqual, 1)
		SkipSo(1, ShouldEqual, 1)
	})

	expectEqual(t, "Begin|A|Skipped|Success|Skipped|Exit|End", myReporter.wholeStory())
}

func TestErrorByManualPanicReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		panic("Gopher alert!")
	})

	expectEqual(t, "Begin|A|Error|Exit|End", myReporter.wholeStory())
}

func TestIterativeConveysReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		for x := 0; x < 3; x++ {
			Convey(strconv.Itoa(x), func() {
				So(x, ShouldEqual, x)
			})
		}
	})

	expectEqual(t, "Begin|A|0|Success|Exit|Exit|A|1|Success|Exit|Exit|A|2|Success|Exit|Exit|End", myReporter.wholeStory())
}

func expectEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected '%v' to be '%v' but it wasn't. See '%s' at line %d.",
			actual, expected, path.Base(file), line)
	}
}

func setupFakeReporter() (*fakeReporter, *fakeGoTest) {
	myReporter := fakeReporter{}
	myReporter.calls = []string{}
	reporter = &myReporter
	runner = execution.NewRunner()
	runner.UpgradeReporter(reporter)
	return &myReporter, &fakeGoTest{}
}

type fakeReporter struct {
	calls []string
}

func (self *fakeReporter) BeginStory(story *reporting.StoryReport) {
	self.calls = append(self.calls, "Begin")
}

func (self *fakeReporter) Enter(scope *reporting.ScopeReport) {
	self.calls = append(self.calls, scope.Title)
}

func (self *fakeReporter) Report(report *reporting.AssertionResult) {
	if report.Error != nil {
		self.calls = append(self.calls, "Error")
	} else if report.Failure != "" {
		message := "Failure"
		if report.Expected != "" || report.Actual != "" {
			message += fmt.Sprintf("(%s/%s)", report.Expected, report.Actual)
		}
		self.calls = append(self.calls, message)
	} else if report.Skipped {
		self.calls = append(self.calls, "Skipped")
	} else {
		self.calls = append(self.calls, "Success")
	}
}

func (self *fakeReporter) Exit() {
	self.calls = append(self.calls, "Exit")
}

func (self *fakeReporter) EndStory() {
	self.calls = append(self.calls, "End")
}

func (self *fakeReporter) wholeStory() string {
	return strings.Join(self.calls, "|")
}
