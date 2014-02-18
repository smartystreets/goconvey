package convey

import "github.com/smartystreets/goconvey/reporting"

type nilReporter struct{}

func (self *nilReporter) BeginStory(story *reporting.StoryReport)  {}
func (self *nilReporter) Enter(scope *reporting.ScopeReport)       {}
func (self *nilReporter) Report(report *reporting.AssertionResult) {}
func (self *nilReporter) Exit()                                    {}
func (self *nilReporter) EndStory()                                {}
func newNilReporter() *nilReporter                                 { return new(nilReporter) }
