package reporting

import (
	"runtime"
	"strings"
)

type Report struct {
	File       string
	Line       int
	Failure    string
	Error      interface{}
	stackTrace string
}

func NewFailureReport(failure string) *Report {
	file, line := caller()
	stack := stackTrace()
	report := Report{file, line, failure, nil, stack}
	return &report
}
func NewErrorReport(err interface{}) *Report {
	file, line := caller()
	stack := fullStackTrace()
	report := Report{file, line, "", err, stack}
	return &report
}
func NewSuccessReport() *Report {
	file, line := caller()
	stack := stackTrace()
	report := Report{file, line, "", nil, stack}
	return &report
}

func caller() (file string, line int) {
	// TODO: what if they have extracted the So() call into a helper method?
	//       (runtime.Caller(3) will not yield the correct stack entry!)
	_, file, line, _ = runtime.Caller(3)
	return
}
func stackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, false)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return filterStack(formatted)
}
func fullStackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, true)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return filterStack(formatted)
}
func filterStack(stack string) string {
	lines := strings.Split(stack, newline)
	filtered := []string{}
	for _, line := range lines {
		if isExternal(line) {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, newline)
}
func isExternal(line string) bool {
	if strings.Contains(line, "goconvey/convey") {
		return false
	} else if strings.Contains(line, "goconvey/execution") {
		return false
	} else if strings.Contains(line, "goconvey/gotest") {
		return false
	} else if strings.Contains(line, "goconvey/printing") {
		return false
	} else if strings.Contains(line, "goconvey/reporting") {
		return false
	}
	return true
}

const newline = "\n"
const success = "✓"
const failure = "✗"
const error_ = "🔥"
const errorTemplate = "* %s \n* Line: %d - %v \n%s"
const failureTemplate = "* %s \n* Line %d: %s"
