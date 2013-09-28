package reporting

import (
	"bytes"
	"encoding/json"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *jsonReporter) BeginStory(test gotest.T) {
	self.story = newScope("TOP")
	self.current = self.story
}

func (self *jsonReporter) Enter(title, id string) {
	if _, found := self.titlesById[id]; found {
		return
	}

	self.titlesById[id] = title
	child := newScope(title)
	self.current.Children = append(self.current.Children, child)
	self.current = child
}

func (self *jsonReporter) Report(r *Report) {
	self.current.Reports = append(self.current.Reports, r)
}

func (self *jsonReporter) Exit() {
	if self.current.parent != nil {
		self.current = self.current.parent
	}
}

func (self *jsonReporter) EndStory() {
	serialized, _ := json.Marshal(self.story)
	var b bytes.Buffer
	json.Indent(&b, serialized, "", " ")
	self.out.Print(b.String())

	self.titlesById = make(map[string]string)
	self.story = nil
	self.current = nil
}

func NewJsonReporter(out *printing.Printer) *jsonReporter {
	self := &jsonReporter{}
	self.titlesById = make(map[string]string)
	self.out = out
	return self
}

type jsonReporter struct {
	out        *printing.Printer
	story      *scope
	current    *scope
	titlesById map[string]string
}

type scope struct {
	Title    string
	Reports  []*Report
	Children []*scope
	parent   *scope
}

func newScope(title string) *scope {
	self := &scope{}
	self.Title = title
	self.Reports = []*Report{}
	self.Children = []*scope{}
	return self
}
