package reporting

import "github.com/smartystreets/goconvey/gotest"
import "github.com/smartystreets/goconvey/printing"

func (self *dot) BeginStory(test gotest.T) {}

func (self *dot) Enter(title, id string) {}

func (self *dot) Report(r *Report) {
	if r.Error != nil {
		self.out.Insert("E")
		self.errors = append(self.errors, r)
	} else if r.Failure != "" {
		self.out.Insert("X")
		self.failures = append(self.failures, r)
	} else {
		self.out.Insert(".")
	}
}

func (self *dot) Exit() {}

func (self *dot) EndStory() {
	self.out.Println("")
	self.showErrors()
	self.showFailures()
}
func (self *dot) showErrors() {
	for i, e := range self.errors {
		if i == 0 {
			self.out.Println("Errors:")
		}
		self.out.Println(errorTemplate, e.File, e.Line, e.Error, e.stackTrace)
	}
}
func (self *dot) showFailures() {
	for i, f := range self.failures {
		if i == 0 {
			self.out.Println("Failures:")
		}
		self.out.Println(failureTemplate, f.File, f.Line, f.Failure)
	}
}

func NewDotReporter(out *printing.Printer) *dot {
	self := dot{}
	self.out = out
	self.errors = []*Report{}
	self.failures = []*Report{}
	return &self
}

type dot struct {
	out      *printing.Printer
	errors   []*Report
	failures []*Report
}
