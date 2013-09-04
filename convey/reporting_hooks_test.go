package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/reporting"
	"path"
	"runtime"
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

func TestErrorByManualPanicReported(t *testing.T) {
	myReporter, test := setupFakeReporter()

	Convey("A", test, func() {
		panic("Gopher alert!")
	})

	expectEqual(t, "Begin|A|Error|Exit|End", myReporter.wholeStory())
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

func (self *fakeReporter) BeginStory(test gotest.T) {
	self.calls = append(self.calls, "Begin")
}

func (self *fakeReporter) Enter(title, id string) {
	self.calls = append(self.calls, title)
}

func (self *fakeReporter) Report(r *reporting.Report) {
	if r.Error != nil {
		self.calls = append(self.calls, "Error")
	} else if r.Failure != "" {
		self.calls = append(self.calls, "Failure")
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
