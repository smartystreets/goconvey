package reporting

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/smartystreets/goconvey/gotest"
)

type FailureView struct {
	Message  string
	Expected string
	Actual   string
}

type AssertionResult struct {
	File       string
	Line       int
	Expected   string
	Actual     string
	Failure    string
	Error      interface{}
	StackTrace string
	Skipped    bool
}

func NewFailureReport(failure string) *AssertionResult {
	report := &AssertionResult{}
	report.File, report.Line = caller()
	report.StackTrace = stackTrace()
	parseFailure(failure, report)
	return report
}
func parseFailure(failure string, report *AssertionResult) {
	view := &FailureView{}
	err := json.Unmarshal([]byte(failure), view)
	if err == nil {
		report.Failure = view.Message
		report.Expected = view.Expected
		report.Actual = view.Actual
	} else {
		report.Failure = failure
	}
}
func NewErrorReport(err interface{}) *AssertionResult {
	report := &AssertionResult{}
	report.File, report.Line = caller()
	report.StackTrace = fullStackTrace()
	report.Error = fmt.Sprintf("%v", err)
	return report
}
func NewSuccessReport() *AssertionResult {
	report := &AssertionResult{}
	report.File, report.Line = caller()
	report.StackTrace = fullStackTrace()
	return report
}
func NewSkipReport() *AssertionResult {
	report := &AssertionResult{}
	report.File, report.Line = caller()
	report.StackTrace = fullStackTrace()
	report.Skipped = true
	return report
}

func caller() (file string, line int) {
	file, line, _ = gotest.ResolveExternalCaller()
	return
}
func stackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, false)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return removeInternalEntries(formatted)
}
func fullStackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, true)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return removeInternalEntries(formatted)
}
func removeInternalEntries(stack string) string {
	lines := strings.Split(stack, newline)
	filtered := []string{}
	for _, line := range lines {
		if !isExternal(line) {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, newline)
}
func isExternal(line string) bool {
	for _, p := range internalPackages {
		if strings.Contains(line, p) {
			return true
		}
	}
	return false
}

// NOTE: any new packages that host goconvey packages will need to be added here!
// An alternative is to scan the goconvey directory and then exclude stuff like
// the examples package but that's nasty too.
var internalPackages = []string{
	"goconvey/assertions",
	"goconvey/convey",
	"goconvey/execution",
	"goconvey/gotest",
	"goconvey/printing",
	"goconvey/reporting",
}
