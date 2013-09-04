package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *statistics) BeginStory(test gotest.T) {
}

func (self *statistics) Enter(title, id string) {}

func (self *statistics) Report(r *Report) {
	self.total++
}

func (self *statistics) Exit() {}

func (self *statistics) EndStory() {
	plural := "s"
	if self.total == 1 {
		plural = ""
	}
	self.out.Println("\n%d assertion%s and counting\n", self.total, plural)
}

func NewStatisticsReporter(out *printing.Printer) *statistics {
	self := statistics{}
	self.out = out
	return &self
}

type statistics struct {
	out   *printing.Printer
	total int
}
