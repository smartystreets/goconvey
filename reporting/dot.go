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
	self.out.Println("Errors:")

	for _, e := range self.errors {
		self.out.Println(errorTemplate, e.File, e.Line, e.Error, e.stackTrace)
	}

	self.out.Println("Failures:")
	for _, f := range self.failures {
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
