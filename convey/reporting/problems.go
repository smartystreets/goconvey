package reporting

import (
	"fmt"
)

type problem struct {
	out      *Printer
	errors   []*AssertionResult
	failures []*AssertionResult
}

func (self *problem) BeginStory()              {}
func (self *problem) Enter(scope *ScopeReport) {}
func (self *problem) Exit()                    {}

func (self *problem) Report(report *AssertionResult) {
	if report.Error != nil {
		self.errors = append(self.errors, report)
	} else if report.Failure != "" {
		self.failures = append(self.failures, report)
	}
}

func (self *problem) Close() {
	self.show(self.showErrors, redColor)
	self.show(self.showFailures, yellowColor)
}
func (self *problem) show(display func(), color string) {
	self.out.Insert(color)
	display()
	self.out.Insert(resetColor)
	self.out.Exit()
}
func (self *problem) showErrors() {
	for i, e := range self.errors {
		if i == 0 {
			self.out.Suite("Errors:")
		}
		self.out.Statement(
			fmt.Sprintf(errorTemplate, e.File, e.Line, e.Error, e.StackTrace))
		self.out.Statement()
	}
}
func (self *problem) showFailures() {
	for i, f := range self.failures {
		if i == 0 {
			self.out.Suite("Failures:")
		}

		self.out.Statement(
			fmt.Sprintf(failureTemplate, f.File, f.Line, f.Failure))
		self.out.Statement()
	}
}

func (self *problem) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewProblemReporter(out *Printer) *problem {
	return &problem{out: out}
}
