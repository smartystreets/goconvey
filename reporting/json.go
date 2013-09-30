package reporting

import (
	"bytes"
	"encoding/json"
	"github.com/smartystreets/goconvey/printing"
)

func (self *jsonReporter) BeginStory(story *StoryReport) {
	top := newScope(story.Name, self.depth, story.File, story.Line)
	self.scopes = append(self.scopes, top)
	self.stack = append(self.stack, top)
}

func (self *jsonReporter) Enter(scope *ScopeReport) {
	self.depth++
	if _, found := self.titlesById[scope.ID]; !found {
		self.registerScope(scope)
	}
}
func (self *jsonReporter) registerScope(scope *ScopeReport) {
	self.titlesById[scope.ID] = scope.ID
	next := newScope(scope.Title, self.depth, scope.File, scope.Line)
	self.scopes = append(self.scopes, next)
	self.stack = append(self.stack, next)
}

func (self *jsonReporter) Report(report *AssertionReport) {
	current := self.stack[len(self.stack)-1]
	current.Reports = append(current.Reports, newAssertion(report))
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
	serialized, _ := json.Marshal(self.scopes)
	var buffer bytes.Buffer
	json.Indent(&buffer, serialized, "", "  ")
	self.out.Print(buffer.String() + ",")
}
func (self *jsonReporter) reset() {
	self.titlesById = make(map[string]string)
	self.scopes = []*scope{}
	self.stack = []*scope{}
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
	scopes     []*scope
	stack      []*scope
	depth      int
}

type scope struct {
	Title   string
	File    string
	Line    int
	Depth   int
	Reports []*assertion
}

func newScope(title string, depth int, file string, line int) *scope {
	self := &scope{}
	self.Title = title
	self.Depth = depth
	self.File = file
	self.Line = line
	self.Reports = []*assertion{}
	return self
}

type assertion struct {
	File    string
	Line    int
	Failure string
	Error   interface{}
	Skipped bool

	// TODO: I'm going to have to parse this turn it into a structure that
	// can accomodate turning the file paths into urls when templated...
	StackTrace string
}

func newAssertion(report *AssertionReport) *assertion {
	self := &assertion{}
	self.File = report.File
	self.Line = report.Line
	self.Failure = report.Failure
	self.Error = report.Error
	self.StackTrace = report.stackTrace
	self.Skipped = report.Skipped
	return self
}
