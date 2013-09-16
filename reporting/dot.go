package reporting

import "fmt"

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *dot) BeginStory(test gotest.T) {}

func (self *dot) Enter(title, id string) {}

func (self *dot) Report(r *Report) {
	if r.Error != nil {
		fmt.Print(redColor)
		self.out.Insert(dotError)
	} else if r.Failure != "" {
		fmt.Print(yellowColor)
		self.out.Insert(dotFailure)
	} else {
		fmt.Print(greenColor)
		self.out.Insert(dotSuccess)
	}
	fmt.Print(resetColor)
}

func (self *dot) Exit() {}

func (self *dot) EndStory() {}

func NewDotReporter(out *printing.Printer) *dot {
	self := dot{}
	self.out = out
	return &self
}

type dot struct {
	out *printing.Printer
}
