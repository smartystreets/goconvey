package reporting

import "github.com/smartystreets/goconvey/gotest"

func (self *gotestReporter) BeginStory(story *StoryReport) {
	self.test = story.Test
}

func (self *gotestReporter) Enter(scope *ScopeReport) {}

func (self *gotestReporter) Report(r *AssertionResult) {
	if !passed(r) {
		self.test.Fail()
	}
}

func (self *gotestReporter) Exit() {}

func (self *gotestReporter) EndStory() {
	self.test = nil
}

func NewGoTestReporter() *gotestReporter {
	self := gotestReporter{}
	return &self
}

type gotestReporter struct {
	test gotest.T
}

func passed(r *AssertionResult) bool {
	return r.Error == nil && r.Failure == ""
}
