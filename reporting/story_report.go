package reporting

import (
	"strings"

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
	name = removePackagePath(name)
	self := &StoryReport{}
	self.Test = test
	self.Name = name
	self.File = file
	self.Line = line
	return self
}

// name comes in looking like "github.com/smartystreets/goconvey/examples.TestName".
// We only want the stuff after the last '.', which is the name of the test function.
func removePackagePath(name string) string {
	parts := strings.Split(name, ".")
	if len(parts) == 1 {
		return name
	}
	return parts[len(parts)-1]
}
