package reporting

import (
	"bytes"
	"encoding/json"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *jsonReporter) BeginStory(test gotest.T) {
	top := newScope("TOP", self.depth) // TODO: this could have the GoTest name.
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
	next := newScope(title, self.depth)
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

func newScope(title string, depth int) *scope {
	self := &scope{}
	self.Depth = depth
	self.Title = title
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
