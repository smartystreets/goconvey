package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
	"time"
)

func (self *statistics) BeginStory(test gotest.T) {

}

func (self *statistics) Enter(title, id string) {

}

func (self *statistics) Report(r *Report) {

}

func (self *statistics) Exit() {

}

func (self *statistics) EndStory() {

}

func NewStatisticsReporter() *statistics {
	self := statistics{}
	self.Reports = make(map[string]*scopeReport)
	return &self
}

type statistics struct {
	Reports map[string]*scopeReport
}

type scopeReport struct {
	name      string
	children  []*scopeReport
	successes []Report
	failures  []Report
	errors    []Report
	duration  time.Duration
}
