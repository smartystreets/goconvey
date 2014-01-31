package reporting

import (
	"fmt"

	"github.com/smartystreets/goconvey/gotest"
)

type ScopeReport struct {
	Title string
	ID    string
	File  string
	Line  int
}

func NewScopeReport(title, name string) *ScopeReport {
	file, line, _ := gotest.ResolveExternalCaller()
	self := &ScopeReport{}
	self.Title = title
	self.ID = fmt.Sprintf("%s|%s", title, name)
	self.File = file
	self.Line = line
	return self
}
