package reporting

import (
	"fmt"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
	"time"
)

func (self *statistics) BeginStory(test gotest.T) {
	self.start = time.Now()
}

func (self *statistics) Enter(title, id string) {}

func (self *statistics) Report(r *Report) {
	if r.Error != nil {
		self.errors++
	} else if r.Failure != "" {
		self.failures++
	} else {
		self.successes++
	}
}

func (self *statistics) Exit() {}

func (self *statistics) EndStory() {
	duration := time.Since(self.start)

	message := fmt.Sprintf("\nPassed: %d", self.successes)
	if self.failures > 0 {
		message += fmt.Sprintf(" | Failed: %d", self.failures)
	}
	if self.errors > 0 {
		message += fmt.Sprintf(" | Errors: %d", self.errors)
	}
	message += fmt.Sprintf(" (in %v)\n", duration)
	self.out.Println(message)
}

func NewStatisticsReporter(out *printing.Printer) *statistics {
	self := statistics{}
	self.out = out
	return &self
}

type statistics struct {
	out       *printing.Printer
	start     time.Time
	successes int
	failures  int
	errors    int
}
