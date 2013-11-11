package execution

import "github.com/smartystreets/goconvey/reporting"

func NewNilReporter() *nilReporter {
	self := nilReporter{}
	return &self
}

func (self *nilReporter) BeginStory(story *reporting.StoryReport)  {}
func (self *nilReporter) Enter(scope *reporting.ScopeReport)       {}
func (self *nilReporter) Report(report *reporting.AssertionResult) {}
func (self *nilReporter) Exit()                                    {}
func (self *nilReporter) EndStory()                                {}

type nilReporter struct{}
