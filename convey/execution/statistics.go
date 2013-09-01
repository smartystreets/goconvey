package execution

import (
	"time"
)

func (self *statistics) Enter(scope string) {

}

func (self *statistics) Report(r *Report) {

}

func (self *statistics) Exit() {

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
