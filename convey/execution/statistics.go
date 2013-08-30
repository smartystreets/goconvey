package execution

import (
	"time"
)

func (self *statistics) Success(scope string) {

}

func (self *statistics) Failure(scope string, problem error) {

}

func (self *statistics) Error(scope string, problem error) {

}

func (self *statistics) End(scope string) {

}

func NewStatisticsReporter() *statistics {
	self := statistics{}
	self.Reports = make(map[string]*report)
	return &self
}

type statistics struct {
	Successes int
	Failures  int
	Errors    int
	Duration  time.Duration
	Reports   map[string]*report
	innert    Reporter
}

type report struct {
	Name       string
	SubReports []*report
	Successes  int
	Failures   []error
	Errors     []error
	Duration   time.Duration
}
