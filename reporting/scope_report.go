package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
)

type ScopeReport struct {
	Title string
	ID    string
	File  string
	Line  int
}

func NewScopeReport(title, id string) *ScopeReport {
	file, line, _ := gotest.ResolveExternalCaller()
	self := &ScopeReport{}
	self.Title = title
	self.ID = id
	self.File = file
	self.Line = line
	return self
}
