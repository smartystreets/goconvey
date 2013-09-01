package convey

import (
	"github.com/smartystreets/goconvey/convey/execution"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestSingleScopeReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		So(1, ShouldEqual, 1)
	})

	expectEqual(t, "A|Success|Exit", reporter.wholeStory())
}

func TestNestedScopeReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		Convey("B", func() {
			So(1, ShouldEqual, 1)
		})
	})

	expectEqual(t, "A|B|Success|Exit|Exit", reporter.wholeStory())
}

func TestFailureReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		So(1, ShouldBeNil)
	})

	expectEqual(t, "A|Failure|Exit", reporter.wholeStory())
}

func TestNestedFailureReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		Convey("B", func() {
			So(2, ShouldBeNil)
		})
	})

	expectEqual(t, "A|B|Failure|Exit|Exit", reporter.wholeStory())
}

func TestSuccessAndFailureReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		So(1, ShouldBeNil)
		So(nil, ShouldBeNil)
	})

	expectEqual(t, "A|Failure|Success|Exit", reporter.wholeStory())
}

func TestErrorByManualPanicReported(t *testing.T) {
	reporter := setupFakeReporter()

	Convey("A", t, func() {
		panic("Gopher alert!")
	})

	expectEqual(t, "A|Error|Exit", reporter.wholeStory())
}

// TODO: test failures, errors, nested failures, nested errors, (ensure cleanup is happening?)

func expectEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected '%v' to be '%v' but it wasn't. See '%s' at line %d.",
			actual, expected, path.Base(file), line)
	}
}

func reportEqual(t *testing.T, expected *execution.Report, actual *execution.Report) {
	if actual.File != expected.File {
		t.Errorf("")
	}
}

func setupFakeReporter() *fakeReporter {
	reporter := fakeReporter{}
	reporter.calls = []string{}
	execution.SpecRunner = execution.NewScopeRunner()
	execution.SpecRunner.UpgradeReporter(&reporter)
	execution.SpecReporter = &reporter
	return &reporter
}

type fakeReporter struct {
	calls []string
}

func (self *fakeReporter) Enter(scope string) {
	self.calls = append(self.calls, scope)
}

func (self *fakeReporter) Report(r *execution.Report) {
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

func (self *fakeReporter) wholeStory() string {
	return strings.Join(self.calls, "|")
}
