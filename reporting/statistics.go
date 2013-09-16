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
}

func (self *statistics) Exit() {}

func (self *statistics) EndStory() {
	plural := "s"
	if self.total == 1 {
		plural = ""
	}
	if self.failing && !self.erroring {
		fmt.Print(yellowColor)
	} else if self.erroring {
		fmt.Print(redColor)
	} else {
		fmt.Print(greenColor)
	}
	self.out.Println("\n%d assertion%s and counting\n", self.total, plural)
	fmt.Print(resetColor)
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
}
