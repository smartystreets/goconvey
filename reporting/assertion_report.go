package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
	"runtime"
	"strings"
)

type AssertionReport struct {
	File       string
	Line       int
	Failure    string
	Error      interface{}
	stackTrace string
	Skipped    bool
}

func NewFailureReport(failure string) *AssertionReport {
	file, line := caller()
	stack := stackTrace()
	report := AssertionReport{file, line, failure, nil, stack, false}
	return &report
}
func NewErrorReport(err interface{}) *AssertionReport {
	file, line := caller()
	stack := fullStackTrace()
	report := AssertionReport{file, line, "", err, stack, false}
	return &report
}
func NewSuccessReport() *AssertionReport {
	file, line := caller()
	stack := stackTrace()
	report := AssertionReport{file, line, "", nil, stack, false}
	return &report
}
func NewSkipReport() *AssertionReport {
	file, line := caller()
	stack := stackTrace()
	report := AssertionReport{file, line, "", nil, stack, true}
	return &report
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
