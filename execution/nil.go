package execution

import "github.com/smartystreets/goconvey/gotest"
import "github.com/smartystreets/goconvey/reporting"

func NewNilReporter() *nilReporter {
	self := nilReporter{}
	return &self
}

func (self *nilReporter) BeginStory(test gotest.T)            {}
func (self *nilReporter) Enter(title, id string)              {}
func (self *nilReporter) Report(r *reporting.AssertionReport) {}
func (self *nilReporter) Exit()                               {}
func (self *nilReporter) EndStory()                           {}

type nilReporter struct{}
