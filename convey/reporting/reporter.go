package reporting

import (
	"time"
)

var DefaultReporter reporter

func init() {
	reporter = NewReporter()
}

type reporter interface {
	Success(scope string)
	Failure(scope string, problem error)
	Error(scope string, problem error)
	End(scope string)
}

func (self *Reporter) Success(scope string) {

}

func (self *Reporter) Failure(scope string, problem error) {

}

func (self *Reporter) Error(scope string, problem error) {

}

func (self *Reporter) End(scope string) {

}

func NewReporter() *Reporter {
	self := Reporter{}
	self.Reports = make(map[string]*Report)
	return &self
}

type Reporter struct {
	Successes int
	Failures  int
	Errors    int
	Duration  time.Duration
	Reports   map[string]*Report
}

type Report struct {
	Name       string
	SubReports []*Report
	Successes  int
	Failures   []error
	Errors     []error
	Duration   time.Duration
}
