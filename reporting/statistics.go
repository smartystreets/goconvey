package reporting

import "fmt"

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *statistics) BeginStory(test gotest.T) {
}

func (self *statistics) Enter(title, id string) {}

func (self *statistics) Report(r *Report) {
	self.total++
	if !self.failing && r.Failure != "" {
		self.failing = true
	}
	if !self.erroring && r.Error != nil {
		self.erroring = true
	}
	if r.Skipped {
		self.skipped += 1
	}
}

func (self *statistics) Exit() {}

func (self *statistics) EndStory() {
	self.reportAssertions()
	self.reportSkippedSections()
	self.completeReport()
	self.completeReport()
}
func (self *statistics) reportAssertions() {
	if self.failing && !self.erroring {
		fmt.Print(yellowColor)
	} else if self.erroring {
		fmt.Print(redColor)
	} else {
		fmt.Print(greenColor)
	}
	self.out.Print("\n%d %s thus far", self.total, plural("assertion", self.total))
}
func (self *statistics) reportSkippedSections() {
	if self.skipped > 0 {
		fmt.Print(yellowColor)
		self.out.Print(" (one or more sections skipped)")
		self.skipped = 0
	}
}
func (self *statistics) completeReport() {
	self.out.Print("\n")
	fmt.Print(resetColor)
}

func plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}

func NewStatisticsReporter(out *printing.Printer) *statistics {
	self := statistics{}
	self.out = out
	return &self
}

type statistics struct {
	out      *printing.Printer
	total    int
	failing  bool
	erroring bool
	skipped  int
}
