package convey

import (
	"github.com/smartystreets/goconvey/convey/reporting"
)

type nilReporter struct{}

func (n *nilReporter) BeginStory(story *reporting.StoryReport)  {}
func (n *nilReporter) Enter(scope *reporting.ScopeReport)       {}
func (n *nilReporter) Report(report *reporting.AssertionResult) {}
func (n *nilReporter) Exit()                                    {}
func (n *nilReporter) EndStory()                                {}
func (n *nilReporter) Write(p []byte) (int, error)              { return len(p), nil }
func newNilReporter() *nilReporter                              { return &nilReporter{} }
