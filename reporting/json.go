package reporting

import (
	"bytes"
	"encoding/json"
	"github.com/smartystreets/goconvey/printing"
	"strings"
)

func (self *jsonReporter) BeginStory(story *StoryReport) {}

func (self *jsonReporter) Enter(scope *ScopeReport) {
	if _, found := self.titlesById[scope.ID]; !found {
		self.registerScope(scope)
	}
	self.depth++
}
func (self *jsonReporter) registerScope(scope *ScopeReport) {
	self.titlesById[scope.ID] = scope.ID
	next := newScopeResult(scope.Title, self.depth, scope.File, scope.Line)
	self.scopes = append(self.scopes, next)
	self.stack = append(self.stack, next)
}

func (self *jsonReporter) Report(report *AssertionReport) {
	current := self.stack[len(self.stack)-1]
	current.Assertions = append(current.Assertions, newAssertionResult(report))
}

func (self *jsonReporter) Exit() {
	self.depth--
	if len(self.stack) > 0 {
		self.stack = self.stack[:len(self.stack)-1]
	}
}

func (self *jsonReporter) EndStory() {
	self.report()
	self.reset()
}
func (self *jsonReporter) report() {
	self.out.Print(OpenJson + "\n")
	scopes := []string{}
	for _, scope := range self.scopes {
		serialized, _ := json.Marshal(scope)
		var buffer bytes.Buffer
		json.Indent(&buffer, serialized, "", "  ")
		scopes = append(scopes, buffer.String())
	}
	self.out.Print(strings.Join(scopes, ",") + ",\n")
	self.out.Print(CloseJson + "\n")
}
func (self *jsonReporter) reset() {
	self.titlesById = make(map[string]string)
	self.scopes = []*ScopeResult{}
	self.stack = []*ScopeResult{}
	self.depth = 0
}

func NewJsonReporter(out *printing.Printer) *jsonReporter {
	self := &jsonReporter{}
	self.out = out
	self.reset()
	return self
}

type jsonReporter struct {
	out        *printing.Printer
	titlesById map[string]string
	scopes     []*ScopeResult
	stack      []*ScopeResult
	depth      int
}

type ScopeResult struct {
	Title      string
	File       string
	Line       int
	Depth      int
	Assertions []AssertionResult
}

func newScopeResult(title string, depth int, file string, line int) *ScopeResult {
	self := &ScopeResult{}
	self.Title = title
	self.Depth = depth
	self.File = file
	self.Line = line
	self.Assertions = []AssertionResult{}
	return self
}

type AssertionResult struct {
	File       string
	Line       int
	Failure    string
	Error      interface{}
	Skipped    bool
	StackTrace string
}

func newAssertionResult(report *AssertionReport) AssertionResult {
	self := AssertionResult{}
	self.File = report.File
	self.Line = report.Line
	self.Failure = report.Failure
	self.Error = report.Error
	self.StackTrace = report.stackTrace
	self.Skipped = report.Skipped
	return self
}

const OpenJson = ">>>>>"  // "⌦"
const CloseJson = "<<<<<" // "⌫"
