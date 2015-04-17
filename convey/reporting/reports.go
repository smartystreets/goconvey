package reporting

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/smartystreets/goconvey/convey/gotest"
)

////////////////// ScopeReport ////////////////////

type ScopeReport struct {
	Title string
	File  string
	Line  int
}

func NewScopeReport(title string) *ScopeReport {
	ret := &ScopeReport{Title: title}
	ret.File, ret.Line = gotest.ResolveExternalCaller()
	return ret
}

////////////////// ScopeResult ////////////////////

type ScopeResult struct {
	Title      string
	File       string
	Line       int
	Depth      int
	Assertions []*AssertionResult
	Output     string
	RandomSeed int64 `json:",omitempty"`
}

////////////////// NestedScopeResult ////////////////////

type NestedScopeResult struct {
	Parent *NestedScopeResult `json:"-"`
	Title  string
	File   string
	Line   int

	// Items is a list of either:
	//   string (Output data)
	//   *NestedScopeResult
	//   *AssertionReport
	Items []interface{}
}

func NewNestedScopeResult(parent *NestedScopeResult, report *ScopeReport) *NestedScopeResult {
	return &NestedScopeResult{
		Parent: parent,
		Title:  report.Title,
		File:   report.File,
		Line:   report.Line,
	}
}

type ScopeExit struct {
	Leaving *NestedScopeResult
}

func (s *NestedScopeResult) Walk(cb func(interface{})) {
	for _, i := range s.Items {
		cb(i)
		if scope, ok := i.(*NestedScopeResult); ok {
			scope.Walk(cb)
			cb(ScopeExit{scope})
		}
	}
}

/////////////////// FailureView ////////////////////////

type FailureView struct {
	Message  string
	Expected string
	Actual   string
}

////////////////////AssertionResult //////////////////////

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

func NewAssertionResult(all bool) *AssertionResult {
	ret := &AssertionResult{}
	ret.File, ret.Line = gotest.ResolveExternalCaller()
	ret.StackTrace = stackTrace(all)
	return ret
}

func NewFailureReport(failure string) *AssertionResult {
	result := NewAssertionResult(false)

	view := &FailureView{}
	err := json.Unmarshal([]byte(failure), view)
	if err == nil {
		result.Failure = view.Message
		result.Expected = view.Expected
		result.Actual = view.Actual
	} else {
		result.Failure = failure
	}
	return result
}
func NewErrorReport(err interface{}) *AssertionResult {
	ret := NewAssertionResult(true)
	ret.Error = fmt.Sprintf("%v", err)
	return ret
}
func NewSuccessReport() *AssertionResult {
	return &AssertionResult{}
}
func NewSkipReport() *AssertionResult {
	ret := NewAssertionResult(false)
	ret.Skipped = true
	return ret
}

func stackTrace(all bool) string {
	buffer := make([]byte, 1024*64)
	n := runtime.Stack(buffer, all)
	lines := strings.Split(string(buffer[:n]), newline)
	filtered := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
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
	"goconvey/reporting",
	"jtolds/gls",
}
