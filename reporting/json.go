package reporting

import (
	"bytes"
	"encoding/json"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *jsonReporter) BeginStory(test gotest.T) {
	file, line, testName := gotest.ResolveExternalCaller()
	top := newScope(testName, self.depth, file, line)
	self.scopes = append(self.scopes, top)
	self.stack = append(self.stack, top)
}

func (self *jsonReporter) Enter(title, id string) {
	self.depth++
	if _, found := self.titlesById[id]; !found {
		self.registerScope(title, id)
	}
}
func (self *jsonReporter) registerScope(title, id string) {
	self.titlesById[id] = title
	file, line, _ := gotest.ResolveExternalCaller()
	// ok, this isn't working--we're getting the wrong line numbers at this
	// point because we are executing from the top-level Convey statement and
	// not the actual Convey statement.
	// because that's where we execute runner.Run() from).
	// I probably need to capture the file and line information when I parse
	// registration in the convey package. That info would be stored on the
	// scope struct and then passed to the reporter during the Enter method
	// so we'd have it here. Although, I hate to modify the signature of the
	// Reporter interface just for the json reporter but I don't see a better
	// way at this point...
	next := newScope(title, self.depth, file, line)
	self.scopes = append(self.scopes, next)
	self.stack = append(self.stack, next)
}

func (self *jsonReporter) Report(r *Report) {
	current := self.stack[len(self.stack)-1]
	current.Reports = append(current.Reports, newJsonReport(r))
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
	Reports []*report
}

func newScope(title string, depth int, file string, line int) *scope {
	self := &scope{}
	self.Title = title
	self.Depth = depth
	self.File = file
	self.Line = line
	self.Reports = []*report{}
	return self
}

type report struct {
	File       string
	Line       int
	Failure    string
	Error      interface{}
	StackTrace string
	Skipped    bool
}

func newJsonReport(r *Report) *report {
	return &report{
		r.File,
		r.Line,
		r.Failure,
		r.Error,
		r.stackTrace,
		r.Skipped,
	}
}
