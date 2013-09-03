package reporting

import "github.com/smartystreets/goconvey/gotest"

func NewNilReporter() *nilReporter {
	self := nilReporter{}
	return &self
}

func (self *nilReporter) BeginStory(test gotest.T) {}
func (self *nilReporter) Enter(title, id string)   {}
func (self *nilReporter) Report(r *Report)         {}
func (self *nilReporter) Exit()                    {}
func (self *nilReporter) EndStory()                {}

type nilReporter struct{}
