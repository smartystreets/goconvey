package convey

import (
	"github.com/smartystreets/goconvey/convey/execution"
	"path"
	"runtime"
	"strings"
	"testing"
)

func TestSingleScopeReported(t *testing.T) {
	reporter, expected := setupFakeReporter()

	Convey("A", t, func() {
		expected.File, expected.Line = currentFileAndNextLine()
		So(1, ShouldEqual, 1)
	})

	expectEqual(t, true, reporter.stories["A"])
	expectEqual(t, 1, len(reporter.stories))

	expectEqual(t, expected, reporter.reports[0])
	expectEqual(t, 1, len(reporter.reports))
}

func TestNestedScopeReported(t *testing.T) {
	reporter, expected := setupFakeReporter()

	Convey("A", t, func() {
		Convey("B", func() {
			expected.File, expected.Line = currentFileAndNextLine()
			So(1, ShouldEqual, 1)
		})
	})

	expectEqual(t, true, reporter.stories["AB"])
	expectEqual(t, 1, len(reporter.stories))

	expectEqual(t, expected, reporter.reports[0])
	expectEqual(t, 1, len(reporter.reports))
}

// TODO: test failures, errors, nested failures, nested errors, (ensure cleanup is happening?)

func expectEqual(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected '%v' to be '%v' but it wasn't. See '%s' at line %d.",
			actual, expected, path.Base(file), line)
	}
}

func expectGreaterThan(t *testing.T, minimum int, higher int) {
	if !(higher > minimum) {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("Expected '%v' to be greater than '%v' but it wasn't. See '%s' at line %d.",
			higher, minimum, path.Base(file), line)
	}
}

func setupFakeReporter() (*fakeReporter, execution.Report) {
	r := fakeReporter{}
	r.scopes = make([]string, 20, 20)
	r.stories = make(map[string]bool)
	r.reports = []execution.Report{}
	execution.SpecRunner = execution.NewScopeRunner()
	execution.SpecRunner.UpgradeReporter(&r)
	execution.SpecReporter = &r
	return &r, execution.Report{}
}

type fakeReporter struct {
	scopeIndex int
	scopes     []string
	reports    []execution.Report
	stories    map[string]bool
}

func (self *fakeReporter) Enter(scope string) {
	self.scopes[self.scopeIndex] = scope
	self.scopeIndex++
}

func (self *fakeReporter) Success(r execution.Report) {
	self.reports = append(self.reports, r)
}

func (self *fakeReporter) Failure(r execution.Report) {
	self.reports = append(self.reports, r)
}

func (self *fakeReporter) Error(r execution.Report) {
	self.reports = append(self.reports, r)
}

func (self *fakeReporter) Exit() {
	self.scopeIndex--
	if self.scopeIndex == 0 {
		self.stories[self.wholeStory()] = true
		self.scopes = make([]string, 20, 20)
	}
}

func (self *fakeReporter) wholeStory() string {
	return strings.Join(self.scopes, "")
}

func currentFileAndNextLine() (string, int) {
	_, file, line, _ := runtime.Caller(1)
	return file, line + 1
}
