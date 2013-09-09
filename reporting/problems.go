package reporting

import "fmt"

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *problem) BeginStory(test gotest.T) {}

func (self *problem) Enter(title, id string) {}

func (self *problem) Report(r *Report) {
	if r.Error != nil {
		self.errors = append(self.errors, r)
	} else if r.Failure != "" {
		self.failures = append(self.failures, r)
	}
}

func (self *problem) Exit() {}

func (self *problem) EndStory() {
	self.out.Println("")
	self.showErrors()
	self.showFailures()
}
func (self *problem) showErrors() {
	fmt.Print(redColor)
	for i, e := range self.errors {
		if i == 0 {
			self.out.Println("\nErrors:\n")
			self.out.Indent()
		}
		self.out.Println(errorTemplate, e.File, e.Line, e.Error, e.stackTrace)
	}
	self.out.Dedent()
	fmt.Print(resetColor)
	self.errors = []*Report{}
}
func (self *problem) showFailures() {
	fmt.Print(redColor)
	for i, f := range self.failures {
		if i == 0 {
			self.out.Println("\nFailures:\n")
			self.out.Indent()
		}
		self.out.Println(failureTemplate, f.File, f.Line, f.Failure)
	}
	self.out.Dedent()
	fmt.Print(resetColor)
	self.failures = []*Report{}
}

func NewProblemReporter(out *printing.Printer) *problem {
	self := problem{}
	self.out = out
	self.errors = []*Report{}
	self.failures = []*Report{}
	return &self
}

type problem struct {
	out      *printing.Printer
	errors   []*Report
	failures []*Report
}
