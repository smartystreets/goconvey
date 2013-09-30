package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
)

type StoryReport struct {
	Test gotest.T
	Name string
	File string
	Line int
}

func NewStoryReport(test gotest.T) *StoryReport {
	file, line, name := gotest.ResolveExternalCaller()
	self := &StoryReport{}
	self.Test = test
	self.Name = name
	self.File = file
	self.Line = line
	return self
}
