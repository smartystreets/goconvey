package reporting

import (
	"fmt"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *story) BeginStory(test gotest.T) {}

func (self *story) Enter(title, id string) {
	self.out.Indent()
	self.currentId = id

	if _, found := self.titlesById[id]; !found {
		self.out.Println("")
		self.out.Print("- " + title)
		self.out.Insert(" ")
		self.titlesById[id] = title
	}
}

func (self *story) Report(r *Report) {
	if r.Error != nil {
		self.reportError(r)
	} else if r.Failure != "" {
		self.reportFailure(r)
	} else {
		self.report(success, "")
	}
}
func (self *story) reportError(r *Report) {
	message := fmt.Sprintf(errorTemplate, r.File, r.Line, r.Error, r.stackTrace)
	self.report(error_, message)
}
func (self *story) reportFailure(r *Report) {
	message := fmt.Sprintf(failureTemplate, r.File, r.Line, r.Failure)
	self.report(failure, message)
}
func (self *story) report(indicator, message string) {
	self.out.Insert(indicator)
	if message == "" {
		return
	}
	self.out.Println("")
	self.out.Indent()
	self.out.Indent()
	self.out.Print(message)
	self.out.Dedent()
	self.out.Dedent()
}

func (self *story) Exit() {
	self.out.Dedent()
}

func (self *story) EndStory() {
	self.currentId = ""
	self.titlesById = make(map[string]string)
	self.out.Println("")
}

func NewStoryReporter(out *printing.Printer) *story {
	self := story{}
	self.out = out
	self.titlesById = make(map[string]string)
	return &self
}

type story struct {
	out        *printing.Printer
	titlesById map[string]string
	currentId  string
}
